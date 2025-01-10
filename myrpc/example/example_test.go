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

package example_test

import (
	"testing"

	"github.com/ivfzhou/study_golang/myrpc/example"
	"github.com/ivfzhou/study_golang/myrpc/server"
)

func TestExample(t *testing.T) {
	svr := server.New()
	err := svr.Register(&example.BookService{})
	if err != nil {
		t.Fatal(err)
	}
	err = svr.ListenAndServe()
	if err != nil {
		t.Fatal(err)
	}
}
