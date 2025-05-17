/*
 * @Author: 5553557@qq.com
 * @Date: 2025-05-10 14:05:37
 * @LastEditors: xxx@xxx.com
 * @LastEditTime: 2025-05-17 18:23:42
 * @FilePath: \domainweb\main.go
 * @Description:
 *
 * Copyright (c) 2022 by 5553557@qq.com, All Rights Reserved.
 */
package main

import (
	"fmt"
	"log"

	"domainweb/cmd"
)

func main() {
	fmt.Println("域名估价系统启动中...")
	if err := cmd.Execute(); err != nil {
		log.Fatalf("程序启动失败: %v", err)
	}
}
