package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"postgres_notify/config"
)

func main() {
	ctx := context.Background()

	// 创建连接池
	pool, err := pgxpool.New(ctx, config.ConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// 获取专用连接用于监听
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Release()

	// 监听频道
	_, err = conn.Exec(ctx, "LISTEN events")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("开始监听通知...")

	for {
		// 等待通知
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("收到通知: Channel=%s Pid=%d Payload=%s\n",
			notification.Channel,
			notification.PID,
			notification.Payload)
	}
}
