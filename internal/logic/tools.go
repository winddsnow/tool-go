package logic

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	v1 "tool-go/api/v1"
	"tool-go/internal/service"

	"github.com/gogf/gf/v2/util/grand"
)

func init() {
	service.RegisterTools(New())
}

type sTools struct{}

func New() *sTools {
	return &sTools{}
}

func (s *sTools) MockData(ctx context.Context, req *v1.MockDataReq) (*v1.MockDataRes, error) {
	typeMap := make(map[string]bool)
	for _, t := range req.Types {
		typeMap[t] = true
	}

	columns := req.Types
	data := make([]map[string]string, 0, req.Count)

	for i := 0; i < req.Count; i++ {
		row := make(map[string]string)
		for _, t := range columns {
			row[t] = s.generateField(t)
		}
		data = append(data, row)
	}

	return &v1.MockDataRes{
		Columns: columns,
		Data:    data,
	}, nil
}

func (s *sTools) generateField(typ string) string {
	switch typ {
	case "name":
		return s.genName()
	case "phone":
		return s.genPhone()
	case "email":
		return s.genEmail()
	case "id_card":
		return s.genIDCard()
	case "passport":
		return s.genPassport()
	case "address":
		return s.genAddress()
	case "ip":
		return s.genIP()
	case "datetime":
		return s.genDateTime()
	case "bank_card":
		return s.genBankCard()
	default:
		return s.genName()
	}
}

var surnames = []string{
	"王", "李", "张", "刘", "陈", "杨", "赵", "黄", "周", "吴",
	"徐", "孙", "胡", "朱", "高", "林", "何", "郭", "马", "罗",
	"梁", "宋", "郑", "谢", "韩", "唐", "冯", "于", "董", "萧",
	"程", "曹", "袁", "邓", "许", "傅", "沈", "曾", "彭", "吕",
	"苏", "卢", "蒋", "蔡", "贾", "丁", "魏", "薛", "叶", "阎",
	"余", "潘", "杜", "戴", "夏", "钟", "汪", "田", "任", "姜",
	"范", "方", "石", "姚", "谭", "廖", "邹", "熊", "金", "陆",
	"郝", "孔", "白", "崔", "康", "毛", "邱", "秦", "江", "史",
}

var givenNames = []string{
	"伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军",
	"洋", "勇", "艳", "杰", "娟", "涛", "明", "超", "秀兰", "霞",
	"平", "刚", "桂英", "文", "华", "飞", "玉兰", "斌", "玲", "国强",
	"志强", "建国", "建华", "志明", "秀梅", "海燕", "红", "丽华", "雪", "思远",
	"梓涵", "子轩", "雨涵", "宇轩", "思琪", "浩宇", "欣怡", "浩然", "诗涵", "子涵",
}

func (s *sTools) genName() string {
	return surnames[rand.Intn(len(surnames))] + givenNames[rand.Intn(len(givenNames))]
}

func (s *sTools) genPhone() string {
	prefixes := []string{"13", "14", "15", "17", "18", "19"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	body := ""
	for i := 0; i < 9; i++ {
		body += strconv.Itoa(rand.Intn(10))
	}
	return prefix + body
}

var emailDomains = []string{
	"qq.com", "163.com", "126.com", "gmail.com", "outlook.com",
	"foxmail.com", "sina.com", "aliyun.com", "yeah.net", "hotmail.com",
	"icloud.com", "proton.me", "live.com", "zoho.com",
}

func (s *sTools) genEmail() string {
	local := ""
	l := rand.Intn(8) + 4
	for i := 0; i < l; i++ {
		local += string(rune('a' + rand.Intn(26)))
	}
	domain := emailDomains[rand.Intn(len(emailDomains))]
	return local + "@" + domain
}

func (s *sTools) genIDCard() string {
	areaCodes := []string{
		"110101", "310101", "440101", "440301", "330101",
		"320101", "510101", "500101", "210101", "350201",
		"420101", "430101", "440401", "440501", "441901",
	}

	area := areaCodes[rand.Intn(len(areaCodes))]
	start := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2004, 12, 31, 0, 0, 0, 0, time.UTC)
	birth := start.Add(time.Duration(rand.Int63n(int64(end.Sub(start)))))
	birthStr := birth.Format("20060102")

	seq := rand.Intn(1000)
	seqStr := fmt.Sprintf("%03d", seq)

	body := area + birthStr + seqStr

	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checksumChars := "10X98765432"
	sum := 0
	for i, c := range body {
		n, _ := strconv.Atoi(string(c))
		sum += n * weights[i]
	}
	check := string(checksumChars[sum%11])

	return body + check
}

func (s *sTools) genPassport() string {
	prefixes := []string{"E", "G", "P", "S", "D"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	body := grand.Digits(8)
	return prefix + body
}

var cityData = []struct {
	province string
	city     string
	district []string
}{
	{"广东省", "广州市", []string{"天河区", "越秀区", "海珠区", "荔湾区", "白云区", "番禺区", "黄埔区", "花都区"}},
	{"广东省", "深圳市", []string{"南山区", "福田区", "罗湖区", "宝安区", "龙岗区", "龙华区", "坪山区", "光明区"}},
	{"广东省", "东莞市", []string{"南城街道", "东城街道", "万江街道", "莞城街道", "长安镇", "虎门镇", "厚街镇", "大朗镇"}},
	{"浙江省", "杭州市", []string{"西湖区", "上城区", "拱墅区", "滨江区", "萧山区", "余杭区", "临平区", "钱塘区"}},
	{"浙江省", "宁波市", []string{"海曙区", "鄞州区", "江北区", "镇海区", "北仑区", "奉化区"}},
	{"上海市", "上海市", []string{"浦东新区", "黄浦区", "徐汇区", "静安区", "长宁区", "普陀区", "虹口区", "杨浦区", "宝山区", "闵行区"}},
	{"北京市", "北京市", []string{"海淀区", "朝阳区", "东城区", "西城区", "丰台区", "通州区", "昌平区", "大兴区"}},
	{"江苏省", "南京市", []string{"玄武区", "秦淮区", "建邺区", "鼓楼区", "栖霞区", "雨花台区", "江宁区", "浦口区"}},
	{"江苏省", "苏州市", []string{"姑苏区", "虎丘区", "吴中区", "相城区", "吴江区", "常熟市", "张家港市", "昆山市"}},
	{"四川省", "成都市", []string{"锦江区", "青羊区", "金牛区", "武侯区", "成华区", "高新区", "天府新区", "龙泉驿区"}},
	{"湖北省", "武汉市", []string{"武昌区", "江汉区", "江岸区", "硚口区", "汉阳区", "洪山区", "青山区", "东西湖区"}},
	{"福建省", "厦门市", []string{"思明区", "湖里区", "集美区", "海沧区", "同安区", "翔安区"}},
	{"山东省", "青岛市", []string{"市南区", "市北区", "李沧区", "崂山区", "城阳区", "黄岛区", "即墨区"}},
	{"湖南省", "长沙市", []string{"岳麓区", "芙蓉区", "天心区", "开福区", "雨花区", "望城区"}},
	{"重庆市", "重庆市", []string{"渝中区", "江北区", "南岸区", "沙坪坝区", "九龙坡区", "大渡口区", "渝北区", "巴南区"}},
}

func (s *sTools) genAddress() string {
	ci := cityData[rand.Intn(len(cityData))]
	dist := ci.district[rand.Intn(len(ci.district))]
	streetNum := rand.Intn(9999) + 1
	return fmt.Sprintf("%s%s%s%d号", ci.province, ci.city, dist, streetNum)
}

func (s *sTools) genIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func (s *sTools) genDateTime() string {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
	t := start.Add(time.Duration(rand.Int63n(int64(end.Sub(start)))))
	return t.Format("2006-01-02 15:04:05")
}

var bankCardPrefixes = []string{
	"622848", "622202", "622188", "621700", "622262",
	"955880", "955888", "622260", "622908", "622155",
	"621226", "622845", "622836", "622161", "621558",
}

func (s *sTools) genBankCard() string {
	prefix := bankCardPrefixes[rand.Intn(len(bankCardPrefixes))]
	body := prefix
	for i := 0; i < 9; i++ {
		body += strconv.Itoa(rand.Intn(10))
	}
	check := s.luhnCheckDigit(body)
	return body + check
}

func (s *sTools) luhnCheckDigit(num string) string {
	sum := 0
	double := false
	for i := len(num) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(num[i]))
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}
	check := (10 - (sum % 10)) % 10
	return strconv.Itoa(check)
}

var _ service.ITools = (*sTools)(nil)
