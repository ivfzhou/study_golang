/*
 * Copyright (c) 2023 ivfzhou
 * myrpc is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package example

import (
	"context"
	"fmt"
)

type BookService struct{}

func (*BookService) Name() string {
	return "BookService"
}

type BuyReq struct {
	Price int
}
type BuyRsp struct {
	Res int
}

func (*BookService) Buy(ctx context.Context, req *BuyReq) (*BuyRsp, error) {
	fmt.Println(req.Price)
	return &BuyRsp{Res: 1}, nil
}
