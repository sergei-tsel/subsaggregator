package service

import (
	"database/sql"
	"reflect"
	"strconv"
	"subsaggregator/internal/model"
	"subsaggregator/internal/repository"
	"subsaggregator/internal/utils"
	"testing"
	"time"
)

func TestGetOneSubscription(t *testing.T) {
	subscriptionRepo := repository.SubscriptionRepoMock{
		Subscriptions: make(map[int]*model.Subscription),
	}

	subs := createTestSubscriptions(subscriptionRepo, 1)

	type args struct {
		subscriptionRepo repository.SubscriptionRepository
		subId            int
	}

	tests := []struct {
		name    string
		args    args
		want    *model.Subscription
		wantErr bool
	}{
		{
			name: "Получение подписки",
			args: args{
				subscriptionRepo: subscriptionRepo,
				subId:            subs[0].Id,
			},
			want:    &subs[0],
			wantErr: false,
		},
		{
			name: "Несуществующая подписка",
			args: args{
				subscriptionRepo: subscriptionRepo,
				subId:            0,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOneSubscription(tt.args.subscriptionRepo, tt.args.subId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetOneSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOneSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListSubscriptions(t *testing.T) {
	subscriptionRepo := repository.SubscriptionRepoMock{
		Subscriptions: make(map[int]*model.Subscription),
		Count:         0,
	}

	subs := createTestSubscriptions(subscriptionRepo, 2)

	type args struct {
		req              ListSubscriptionsRequest
		subscriptionRepo repository.SubscriptionRepository
	}

	tests := []struct {
		name    string
		args    args
		want    []model.Subscription
		wantErr bool
	}{
		{
			name: "Получение всех подписок за выбранный период",
			args: args{
				req: ListSubscriptionsRequest{
					ServiceName: "",
					UserId:      "",
					StartDate:   *subs[0].StartDate,
					EndDate:     *subs[0].EndDate,
				},
				subscriptionRepo: subscriptionRepo,
			},
			want:    []model.Subscription{subs[0], subs[1]},
			wantErr: false,
		},
		{
			name: "Получение отфильтрованных подписок",
			args: args{
				req: ListSubscriptionsRequest{
					ServiceName: subs[0].ServiceName,
					UserId:      subs[0].UserId,
					StartDate:   *subs[0].StartDate,
					EndDate:     *subs[0].EndDate,
				},
				subscriptionRepo: subscriptionRepo,
			},
			want:    []model.Subscription{subs[0]},
			wantErr: false,
		},
		{
			name: "Получение пустого списка подписок",
			args: args{
				req: ListSubscriptionsRequest{
					ServiceName: "Тестовый сервис 3",
					UserId:      "Тестовый UUID 3",
					StartDate:   *subs[0].StartDate,
					EndDate:     *subs[0].EndDate,
				},
				subscriptionRepo: subscriptionRepo,
			},
			want:    []model.Subscription{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListSubscriptions(tt.args.req, tt.args.subscriptionRepo)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListSubscriptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumSubscriptionsPrices(t *testing.T) {
	subscriptionRepo := repository.SubscriptionRepoMock{
		Subscriptions: make(map[int]*model.Subscription),
		Count:         0,
	}

	subs := createTestSubscriptions(subscriptionRepo, 2)

	type args struct {
		req              SumSubscriptionsPricesRequest
		subscriptionRepo repository.SubscriptionRepository
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Получение суммарной стоимости всех подписок",
			args: args{
				req: SumSubscriptionsPricesRequest{
					ServiceName: "",
					UserId:      "",
					StartDate:   *subs[0].StartDate,
					EndDate:     *subs[0].EndDate,
				},
				subscriptionRepo: subscriptionRepo,
			},
			want:    subs[0].Price + subs[1].Price,
			wantErr: false,
		},
		{
			name: "Получение суммарной стоимости отфильтрованных подписок",
			args: args{
				req: SumSubscriptionsPricesRequest{
					ServiceName: subs[0].ServiceName,
					UserId:      subs[0].UserId,
					StartDate:   *subs[0].StartDate,
					EndDate:     *subs[0].EndDate,
				},
				subscriptionRepo: subscriptionRepo,
			},
			want:    subs[0].Price,
			wantErr: false,
		},
		{
			name: "Отфильтрован пустой список записей о подписках",
			args: args{
				req: SumSubscriptionsPricesRequest{
					ServiceName: "Тестовый сервис 3",
					UserId:      "Тестовый UUID 3",
					StartDate:   *subs[0].StartDate,
					EndDate:     *subs[0].EndDate,
				},
				subscriptionRepo: subscriptionRepo,
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SumSubscriptionsPrices(tt.args.req, tt.args.subscriptionRepo)

			if (err != nil) != tt.wantErr {
				t.Errorf("SumSubscriptionsPrices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if *got != tt.want {
				t.Errorf("SumSubscriptionsPrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateSubscription(t *testing.T) {
	subscriptionRepo := repository.SubscriptionRepoMock{
		Subscriptions: make(map[int]*model.Subscription),
	}

	req := CreateSubscriptionRequest{
		ServiceName: "Тестовый сервис",
		Price:       100,
		UserId:      "Тестовый UUID",
		StartDate: &utils.Date{NullTime: sql.NullTime{
			Time:  time.Now().AddDate(0, -6, 0),
			Valid: true,
		}},
		EndDate: &utils.Date{NullTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}},
	}

	type args struct {
		req  CreateSubscriptionRequest
		repo repository.SubscriptionRepository
	}

	tests := []struct {
		name    string
		args    args
		want    *model.Subscription
		wantErr bool
	}{
		{
			name: "Создание подписки",
			args: args{
				req:  req,
				repo: subscriptionRepo,
			},
			want: &model.Subscription{
				ServiceName: req.ServiceName,
				Price:       req.Price,
				UserId:      req.UserId,
				StartDate:   req.StartDate,
				EndDate:     req.EndDate,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSubscription(tt.args.req, tt.args.repo)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.ServiceName, tt.want.ServiceName) ||
				!reflect.DeepEqual(got.UserId, tt.want.UserId) {
				t.Errorf("CreateSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateSubscription(t *testing.T) {
	subscriptionRepo := repository.SubscriptionRepoMock{
		Subscriptions: make(map[int]*model.Subscription),
	}

	subs := createTestSubscriptions(subscriptionRepo, 1)

	req := UpdateSubscriptionRequest{
		Price: 200,
		StartDate: &utils.Date{NullTime: sql.NullTime{
			Time:  time.Now().AddDate(0, -6, 0),
			Valid: true,
		}},
		EndDate: &utils.Date{NullTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}},
	}

	type args struct {
		req    UpdateSubscriptionRequest
		repo   repository.SubscriptionRepository
		subsId int
	}

	tests := []struct {
		name    string
		args    args
		want    *model.Subscription
		wantErr bool
	}{
		{
			name: "Изменение подписки",
			args: args{
				req:    req,
				repo:   subscriptionRepo,
				subsId: subs[0].Id,
			},
			want: &model.Subscription{
				ServiceName: subs[0].ServiceName,
				Price:       req.Price,
				UserId:      subs[0].UserId,
				StartDate:   req.StartDate,
				EndDate:     req.EndDate,
			},
			wantErr: false,
		},
		{
			name: "Несуществующая подписка",
			args: args{
				req: UpdateSubscriptionRequest{
					Price: 0,
					StartDate: &utils.Date{NullTime: sql.NullTime{
						Time:  time.Now().AddDate(0, -6, 0),
						Valid: true,
					}},
					EndDate: &utils.Date{NullTime: sql.NullTime{
						Time:  time.Now(),
						Valid: true,
					}},
				},
				repo:   subscriptionRepo,
				subsId: 0,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateSubscription(tt.args.req, tt.args.repo, tt.args.subsId)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil && tt.want != nil) &&
				(!reflect.DeepEqual(got.Price, tt.want.Price) ||
					!reflect.DeepEqual(got.StartDate, tt.want.StartDate) ||
					!reflect.DeepEqual(got.EndDate, tt.want.EndDate)) {
				t.Errorf("UpdateSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteSubscription(t *testing.T) {
	subscriptionRepo := repository.SubscriptionRepoMock{
		Subscriptions: make(map[int]*model.Subscription),
	}

	subs := createTestSubscriptions(subscriptionRepo, 1)

	type args struct {
		repo  repository.SubscriptionRepository
		subId int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Удаление подписки",
			args: args{
				repo:  subscriptionRepo,
				subId: subs[0].Id,
			},
			wantErr: false,
		},
		{
			name: "Несуществующая подписка",
			args: args{
				repo:  subscriptionRepo,
				subId: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteSubscription(tt.args.repo, tt.args.subId)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func createTestSubscriptions(repo repository.SubscriptionRepoMock, count int) []model.Subscription {
	subs := []model.Subscription{}

	for i := 0; i < count; i++ {
		createReq := CreateSubscriptionRequest{
			ServiceName: "Тестовый сервис " + strconv.Itoa(i+1),
			Price:       100 + (100 * i),
			UserId:      "Тестовый UUID " + strconv.Itoa(i+1),
			StartDate: &utils.Date{NullTime: sql.NullTime{
				Time:  time.Now().AddDate(0, -6, 0),
				Valid: true,
			}},
			EndDate: &utils.Date{NullTime: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}},
		}

		sub, _ := CreateSubscription(
			createReq,
			repo,
		)

		repo.Count++

		subs = append(subs, *sub)
	}

	return subs
}
