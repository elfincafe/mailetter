package mailetter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewMail(t *testing.T) {
	from, _ := NewAddress("mailetter@example.com", "MaiLetter")
	m := NewMail(from)
	typ := reflect.TypeOf(m)
	expect := reflect.TypeOf((*Mail)(nil)).String()
	if typ.String() != expect {
		t.Errorf("[Case%d] Type: %s != %s", 0, typ.String(), expect)
	}
}

func TestMailHeader(t *testing.T) {
	type tcase struct {
		key string
		val string
	}
	cases := []tcase{}
	cases = append(cases, tcase{"Subject", "Mail Subject"})
	cases = append(cases, tcase{"X-Mailer", "MaiLetter Client"})
	cases = append(cases, tcase{"Message-ID", "<1234567890ABCDEFGHIJKLMN@example.com>"})
	cases = append(cases, tcase{"X-Mailer", "MaiLetter Client"})
	from, _ := NewAddress("mailetter@example.com", "MaiLetter")
	m := NewMail(from)
	for k, v := range cases {
		m.Header(v.key, v.val)
		flg := false
		for k, _ := range m.headers {
			label := strings.ToLower(k)
			if label == strings.ToLower(k) {
				flg = true
				break
			}
		}
		if !flg {
			// t.Errorf("[Case%d] Header: %s != %s; Value %s != %s", k)
			t.Errorf("[Case%d]", k)
		}
	}
}

func TestMailTo(t *testing.T) {

}

func TestMailCc(t *testing.T) {

}

func TestMailBcc(t *testing.T) {

}

func TestMailFrom(t *testing.T) {

}

func TestMailSubject(t *testing.T) {

}

func TestMailBody(t *testing.T) {
	type tcase struct {
		body string
	}
	cases := []tcase{}
	cases = append(cases, tcase{body: "TestBody"})

	from, _ := NewAddress("mailetter@example.com", "MaiLetter")
	mail := NewMail(from)
	for k, v := range cases {
		mail.Body(v.body)
		if mail.body != v.body {
			t.Error(fmt.Sprintf("[Case%d] Body: %s (%s)", k, mail.body, v.body))
		}
	}
}

func TestMailSet(t *testing.T) {

}

func TestMailString(t *testing.T) {

	from, _ := NewAddress("from@example.com", "テスト")
	m := NewMail(from)
	for _, v := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		to, _ := NewAddress(fmt.Sprintf("to+%d@example.com", v), "受信者1")
		m.To(to)
	}
	m.Subject("貨表示を円貨にした場合、平均取得価額は国内約定日の10時30分までは参考レートに為替掛目1％を加えて計算している")
	m.Body("通貨表示を円貨とした場合の時価評価額・評価損益は、現在の参考為替レートを利用して円換算額を算出しております。\r\n従って、実際の円貨決済による売却時の円換算レート(TTB)、\r\nまたそれにより算出される売却価額 (売却時の費用を考慮しない) とは異なります。")
	t.Errorf("%s", m.String())
}
