package postgres

import (
	"context"
	"reflect"
	"testing"

	userModel "cpolzer.de/m/v2/pkg/storage/user"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDatabase() (testcontainers.Container, userModel.Repository, error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, nil, err
	}
	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, nil, err
	}

	repository, err := New(context.TODO(), host, port.Port())
	if err != nil {
		return nil, nil, err
	}

	return dbContainer, repository, nil
}

func Test_pg_Create(t *testing.T) {

	dbContainer, repository, err := SetupTestDatabase()
	if err != nil {

	}
	defer dbContainer.Terminate(context.Background())

	type args struct {
		user userModel.UserDto
	}
	tests := []struct {
		name    string
		args    args
		want    *userModel.UserDto
		wantErr bool
	}{
		{
			name: "save simple user expect response",
			args: args{
				user: userModel.UserDto{
					FirstName: "UserFirstName",
					LastName:  "UserLastName",
					Address:   userModel.AddressDto{},
				},
			},
			want: &userModel.UserDto{
				FirstName: "UserFirstName",
				LastName:  "UserLastName",
				Address:   userModel.AddressDto{},
			},
			wantErr: false,
		},
		{
			name: "save user with address expect response",
			args: args{
				user: userModel.UserDto{
					FirstName: "UserFirstName",
					LastName:  "UserLastName",
					Address: userModel.AddressDto{
						Street:     "Street",
						PostalCode: "123456789",
					},
				},
			},
			want: &userModel.UserDto{
				FirstName: "UserFirstName",
				LastName:  "UserLastName",
				Address: userModel.AddressDto{
					Street:     "Street",
					PostalCode: "123456789",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := repository.Create(context.TODO(), &tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.ID = got.ID
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pg_Get(t *testing.T) {
	dbContainer, repository, err := SetupTestDatabase()
	if err != nil {

	}
	defer dbContainer.Terminate(context.Background())

	type args struct {
		userId        uint
		prefillEntity *userModel.UserDto
	}
	tests := []struct {
		name    string
		args    args
		want    *userModel.UserDto
		wantErr bool
	}{
		{
			name: "get user returns expected",
			args: args{
				userId: 1,
				prefillEntity: &userModel.UserDto{
					FirstName: "UserFirstName",
					LastName:  "UserLastName",
					Address: userModel.AddressDto{
						Street:     "Street",
						PostalCode: "Postal",
					},
				},
			},
			want: &userModel.UserDto{
				FirstName: "UserFirstName",
				LastName:  "UserLastName",
				Address: userModel.AddressDto{
					Street:     "Street",
					PostalCode: "Postal",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			created, err2 := repository.Create(context.TODO(), tt.args.prefillEntity)
			if err2 != nil {
				t.Errorf("Problem creating prefillEntity")
			}
			got, err := repository.Get(context.TODO(), tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.ID = created.ID
			tt.want.Address.ID = created.Address.ID
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pg_Update(t *testing.T) {

	dbContainer, repository, err := SetupTestDatabase()
	if err != nil {

	}
	defer dbContainer.Terminate(context.Background())

	type args struct {
		prefillEntity *userModel.UserDto
		updateEntity  *userModel.UserDto
	}
	tests := []struct {
		name    string
		args    args
		want    *userModel.UserDto
		wantErr bool
	}{
		{
			name: "test update",
			args: args{
				prefillEntity: &userModel.UserDto{
					FirstName: "UserFirstName",
					LastName:  "UserLastName",
					Address: userModel.AddressDto{
						Street:     "Street",
						PostalCode: "Postal",
					},
				},
				updateEntity: &userModel.UserDto{
					FirstName: "UPDATED-UserFirstName",
					LastName:  "UPDATED-UserLastName",
					Address: userModel.AddressDto{
						Street:     "UPDATED-Street",
						PostalCode: "UPDATED-Postal",
					},
				},
			},
			want: &userModel.UserDto{
				FirstName: "UPDATED-UserFirstName",
				LastName:  "UPDATED-UserLastName",
				Address: userModel.AddressDto{
					Street:     "UPDATED-Street",
					PostalCode: "UPDATED-Postal",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err2 := repository.Create(context.TODO(), tt.args.prefillEntity)
			if err2 != nil {
				t.Errorf("Problem creating prefillEntity")
			}
			if diff := cmp.Diff(tt.args.prefillEntity, created); diff != "" {
				t.Logf("Diff (-prefillEntity +created):\n%s", diff)
			}

			tt.args.updateEntity.ID = created.ID
			tt.args.updateEntity.Address.ID = created.Address.ID

			got, err := repository.Update(context.TODO(), tt.args.updateEntity)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
