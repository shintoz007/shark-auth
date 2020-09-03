package refreshtoken

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRefreshTokenRepository struct {
	db *sqlx.DB
}

func NewUserRefreshTokenRepository(db *sqlx.DB) UserRefreshTokenRepository {
	return UserRefreshTokenRepository{db: db}
}

func (r *UserRefreshTokenRepository) Create(refreshToken UserRefreshToken) error {
	_, err := r.db.NamedExec("insert into refresh_tokens (refresh_token, user_id, session_id, expires_at, created_at)"+
		" values (:refresh_token, :user_id, :session_id, :expires_at, :created_at)", refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRefreshTokenRepository) Get(refreshToken string) (UserRefreshToken, error) {
	row := r.db.QueryRowx("select * from refresh_tokens where refresh_token=$1 and deleted=false", refreshToken)
	if row.Err() != nil {
		logrus.WithError(row.Err()).
			Error("error while querying DB")
		return UserRefreshToken{}, row.Err()
	}

	var userRefreshToken UserRefreshToken
	err := row.StructScan(&userRefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserRefreshToken{}, nil
		}
		return UserRefreshToken{}, err
	}
	return userRefreshToken, nil
}

func (r *UserRefreshTokenRepository) RemoveBySessionID(sessionID string) error {
	_, err := r.db.Exec("update refresh_tokens set deleted=true where session_id=$1", sessionID)
	if err != nil {
		return err
	}
	return nil
}
