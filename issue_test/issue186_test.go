/*
* Copyright 2021 ByteDance Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package issue_test

import (
    `math/rand`
    `sync`
    `sync/atomic`
    `testing`
    `time`

    `github.com/bytedance/sonic`
)

var target atomic.Value

type GlobalConfig []Conf

type Conf struct {
    A  string             `json:"A"`
    B  SubConf            `json:"B"`
    C []string            `json:"C"`
}

type SubConf struct {
    Slice        []int64 `json:"Slice"`
    Map          map[int64]bool `json:"-"`
}

func IntSlide2Map(l []int64) map[int64]bool {
    m := make(map[int64]bool)
    for _, item := range l {
        m[item] = true
    }
    return m
}

func Reload(t *testing.T, rawData []byte, target *atomic.Value) {
    var tmp GlobalConfig
    err := sonic.Unmarshal((rawData), &tmp) // better use sonic.UnmarshalString()!
    if err != nil {
        t.Fatalf("failed to unmarshal json, raw data: %v, err: %v", rawData, err)
    }
    for index, conf := range tmp {
        tmp[index].B.Map = IntSlide2Map(conf.B.Slice)
        }
    target.Store(tmp)
}

func TestIssue186(t *testing.T) {
    // t.Parallel()
    var data = []byte(`[{"A":"xxx","B":{"Slice":[111]}},{"A":"yyy","B":{"Slice":[222]},"C":["extra"]},{"A":"zzz","B":{"Slice":[333]},"C":["extra"]}]`)

    for k:=0; k<100; k++ {
        wg := sync.WaitGroup{}
        for i:=0; i<10000; i++ {
            wg.Add(1)
            go func(){
                defer wg.Done()
                time.Sleep(time.Duration(rand.Intn(1000)+1000))
                Reload(t, data, &target)
                        }()
        }
        wg.Wait()
    }
}