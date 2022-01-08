package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDetail2(t *testing.T) {

	err := getDetailInfo2("https://fendou.duoduozy.com/m3u8/202201/4/9aefd6f2c21b618657ffc402d10068f23592ed1f.m3u8?st=bRyYq1p_EkTlTcWXwc46pg&e=1641609068")
	assert.Nil(t, err)
}
