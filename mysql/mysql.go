package mysql

import (
	"context"

	gormHelper "github.com/imkouga/kratos-pkg/contrib/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Transactor struct {
	db *gorm.DB
}

type MySQL struct {
	db *gorm.DB
}

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type Option interface {
	GetDebug() bool
	GetDriver() string
	GetSource() string
}

// 用来承载事务的上下文
type contextTxKey struct{}

// NewMySQL .
func NewMySQL(opt Option, logger log.Logger) (*MySQL, func(), error) {
	clog := log.NewHelper(logger)

	db, err := gormHelper.NewWithOptions(opt,
		gormHelper.WithLogger(clog, false),
		gormHelper.WithTracing(),
	)
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		clog.Info("closing the data resources")
	}
	return &MySQL{
		db: db,
	}, cleanup, nil
}

// ExecTx gorm Transaction
func (t *Transactor) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// ExecTx gorm Transaction
func (d *MySQL) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// DB 根据此方法来判断当前的 db 是不是使用 事务的 DB
func (d *MySQL) DB(ctx context.Context) *gorm.DB {

	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if !ok {
		return d.db
	}
	return tx
}
