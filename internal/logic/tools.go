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

// init 包初始化函数，在程序启动时自动注册 Tools 业务逻辑到 service 层。
func init() {
	service.RegisterTools(New())
}

// sTools 是工具类业务逻辑的实现结构体。
// 注意这个结构体没有像 sUser/sRole 那样实现 service 包中的接口，
// 而是直接在 logic 包内部定义的私有类型，通过编译时检查确保它满足 ITools 接口。
type sTools struct{}

// New 构造函数，返回 *sTools 实例。
// 注意与 sUser/sRole 不同：这里返回的是具体类型指针而非接口，
// 因为外部调用方可以直接使用 *sTools 调用未在接口中声明的方法。
func New() *sTools {
	return &sTools{}
}

// MockData 生成模拟数据，支持姓名、手机号、邮箱、身份证、护照、地址、IP、日期时间、银行卡号等类型。
// 根据请求中指定的 Types 列表，逐行逐字段生成随机数据。
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

// generateField 根据字段类型分派到对应的生成函数。
// 使用 switch-case 模式替代 if-else，更清晰且性能更好。
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

// surnames 常见中文姓氏列表，用于随机生成姓名。
// Go 中 var 声明的包级变量在程序启动时初始化，生命周期为整个进程。
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

// givenNames 常见中文名字列表，包含传统名和现代名，用于随机生成姓名。
var givenNames = []string{
	"伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军",
	"洋", "勇", "艳", "杰", "娟", "涛", "明", "超", "秀兰", "霞",
	"平", "刚", "桂英", "文", "华", "飞", "玉兰", "斌", "玲", "国强",
	"志强", "建国", "建华", "志明", "秀梅", "海燕", "红", "丽华", "雪", "思远",
	"梓涵", "子轩", "雨涵", "宇轩", "思琪", "浩宇", "欣怡", "浩然", "诗涵", "子涵",
}

// genName 随机生成中文姓名："姓 + 名"。
// 使用标准库 math/rand 生成随机索引，rand.Intn(n) 返回 [0, n) 的随机整数。
// 注意：math/rand 是伪随机数生成器，适合模拟数据生成场景。
// 相比之下，gitee.com/gogf/gf/v2/util/grand 提供了更安全的随机数（基于 crypto/rand）。
// 本项目同时使用了 rand（简单随机）和 grand（密码学级随机），根据安全性需求选择。
func (s *sTools) genName() string {
	return surnames[rand.Intn(len(surnames))] + givenNames[rand.Intn(len(givenNames))]
}

// genPhone 随机生成中国大陆手机号码：前 3 位为运营商号段 + 9 位随机数字。
// 当前中国手机号规则：11 位数字，以 1 开头，第二位为 3-9。
func (s *sTools) genPhone() string {
	prefixes := []string{"13", "14", "15", "17", "18", "19"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	body := ""
	for i := 0; i < 9; i++ {
		body += strconv.Itoa(rand.Intn(10))
	}
	return prefix + body
}

// emailDomains 常见邮箱域名列表，用于随机生成电子邮件地址。
var emailDomains = []string{
	"qq.com", "163.com", "126.com", "gmail.com", "outlook.com",
	"foxmail.com", "sina.com", "aliyun.com", "yeah.net", "hotmail.com",
	"icloud.com", "proton.me", "live.com", "zoho.com",
}

// genEmail 随机生成电子邮件地址：随机 4-12 位小写字母 @ 随机邮箱域名。
func (s *sTools) genEmail() string {
	local := ""
	l := rand.Intn(8) + 4
	for i := 0; i < l; i++ {
		local += string(rune('a' + rand.Intn(26)))
	}
	domain := emailDomains[rand.Intn(len(emailDomains))]
	return local + "@" + domain
}

// genIDCard 随机生成合法的 18 位中国大陆公民身份证号码。
//
// 身份证号码结构（GB 11643-1999）：
//   1-6 位   → 地址码（省、市、区县）
//   7-14 位  → 出生日期（YYYYMMDD）
//   15-17 位 → 顺序码（同一地区同一天出生的人的顺序编号，奇数为男、偶数为女）
//   18 位    → 校验码（通过前 17 位按加权因子计算得出）
//
// 校验算法（ISO 7064:1983, MOD 11-2）：
//   1. 前 17 位每位乘以对应的加权因子：7,9,10,5,8,4,2,1,6,3,7,9,10,5,8,4,2
//   2. 加权求和，取和模 11
//   3. 根据余数查校验字符表："10X98765432"
//      - 余数 0 → '1'，余数 1 → '0'，余数 2 → 'X'，余数 3 → '9'，...
//   X 代表罗马数字 10，所以有些身份证号以 X 结尾。
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

// genPassport 随机生成护照号码。
// 护照号格式：1 位字母前缀（E/G/P/S/D）+ 8 位数字。
// 这里使用了 grand.Digits(8) 而非 rand，因为 grand 生成的随机数更均匀，
// 适合编号类场景。grand.Digits(n) 生成 n 位纯数字字符串。
func (s *sTools) genPassport() string {
	prefixes := []string{"E", "G", "P", "S", "D"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	body := grand.Digits(8)
	return prefix + body
}

// cityData 城市数据列表，包含省份、城市、区/县/街道，用于随机生成地址。
// 覆盖全国多个主要城市，每个城市包含多个子区域，使生成的地址更真实。
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

// genAddress 随机生成中文地址："省份 + 城市 + 区/县 + 随机门牌号 + 号"。
func (s *sTools) genAddress() string {
	ci := cityData[rand.Intn(len(cityData))]
	dist := ci.district[rand.Intn(len(ci.district))]
	streetNum := rand.Intn(9999) + 1
	return fmt.Sprintf("%s%s%s%d号", ci.province, ci.city, dist, streetNum)
}

// genIP 随机生成 IPv4 地址（点分十进制格式）。
// 注意：生成的 IP 可能属于保留地址段（如 0.x.x.x、127.x.x.x、私有地址等），
// 模拟数据场景下通常无关紧要。
func (s *sTools) genIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

// genDateTime 在 2020-01-01 至 2026-12-31 范围内随机生成日期时间字符串。
// Go 时间格式化使用参考时间：2006-01-02 15:04:05 是固定的格式标记。
func (s *sTools) genDateTime() string {
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
	t := start.Add(time.Duration(rand.Int63n(int64(end.Sub(start)))))
	return t.Format("2006-01-02 15:04:05")
}

// bankCardPrefixes 常见银行 BIN 号（发卡行识别码）前 6 位。
// 中国银联卡以 62 开头，其他银行各有特定前缀。
var bankCardPrefixes = []string{
	"622848", "622202", "622188", "621700", "622262",
	"955880", "955888", "622260", "622908", "622155",
	"621226", "622845", "622836", "622161", "621558",
}

// genBankCard 随机生成银行卡号（含 Luhn 校验位）。
// 银行卡号结构：发卡行 BIN(6位) + 账户号(9位) + Luhn校验位(1位)。
// 总计 16 位，符合中国银联标准借记卡号长度。
func (s *sTools) genBankCard() string {
	prefix := bankCardPrefixes[rand.Intn(len(bankCardPrefixes))]
	body := prefix
	for i := 0; i < 9; i++ {
		body += strconv.Itoa(rand.Intn(10))
	}
	check := s.luhnCheckDigit(body)
	return body + check
}

// luhnCheckDigit 计算 Luhn 校验位（ISO/IEC 7812-1 标准）。
// Luhn 算法广泛用于银行卡号、信用卡号、IMEI 号等场景。
// 算法步骤：
//   1. 从右向左遍历，每隔一位数字乘以 2。
//   2. 如果乘以 2 后的结果大于 9，则减去 9（或将其各位数字相加）。
//   3. 将所有数字求和。
//   4. 校验位 = (10 - 和 % 10) % 10。
// 此函数接收不含校验位的前缀数字，返回计算出的校验位字符串。
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

// 编译时检查：确保 *sTools 实现了 service.ITools 接口。
// Go 语言中，var _ Interface = (*T)(nil) 是一个惯用技巧：
//   将 nil 指针 *sTools 赋值给 ITools 接口变量，如果 *sTools 没有完整实现
//   ITools 接口，编译会报错。这种检查发生在编译期，零运行时开销。
// 这样可以在代码编辑阶段就发现接口实现不匹配的问题，而非等到运行时才暴露。
var _ service.ITools = (*sTools)(nil)
