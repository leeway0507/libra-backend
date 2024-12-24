package scrap

import "libra-backend/model"

func GetInstance(libCode string) func(isbn string) model.LibScrap {
	libInfo, isExist := LibCodeMap[libCode]
	if !isExist {
		return nil
	}
	inst, isExist := InstanceMap[libInfo.District]
	if !isExist {
		return nil
	}
	return func(isbn string) model.LibScrap {
		return inst(isbn, libInfo.District, libInfo.LibName)
	}
}

type ScrapInstance = func(isbn, district, libname string) model.LibScrap

var InstanceMap = map[string]ScrapInstance{
	"교육청": NewEduction,
	"양천구": NewYangcheon,
}

type lib struct {
	District string
	LibName  string
}

var LibCodeMap = map[string]lib{
	"111003": {
		District: "교육청",
		LibName:  "서울특별시교육청강남도서관",
	},
	"111004": {
		District: "교육청",
		LibName:  "서울특별시교육청강동도서관",
	},
	"111005": {
		District: "교육청",
		LibName:  "서울특별시교육청강서도서관",
	},
	"111006": {
		District: "교육청",
		LibName:  "서울특별시교육청개포도서관",
	},
	"111007": {
		District: "교육청",
		LibName:  "서울특별시교육청고덕평생학습관",
	},
	"111008": {
		District: "교육청",
		LibName:  "서울특별시교육청고척도서관",
	},
	"111009": {
		District: "교육청",
		LibName:  "서울특별시교육청구로도서관",
	},
	"111010": {
		District: "교육청",
		LibName:  "서울특별시교육청남산도서관",
	},
	"111011": {
		District: "교육청",
		LibName:  "서울특별시교육청도봉도서관",
	},
	"111012": {
		District: "교육청",
		LibName:  "서울특별시교육청동대문도서관",
	},
	"111013": {
		District: "교육청",
		LibName:  "서울특별시교육청동작도서관",
	},
	"111014": {
		District: "교육청",
		LibName:  "서울특별시교육청마포평생학습관",
	},
	"111015": {
		District: "교육청",
		LibName:  "서울특별시교육청양천도서관",
	},
	"111016": {
		District: "교육청",
		LibName:  "서울특별시교육청서대문도서관",
	},
	"111017": {
		District: "교육청",
		LibName:  "서울특별시교육청서울시립어린이도서관",
	},
	"111018": {
		District: "교육청",
		LibName:  "서울특별시교육청영등포평생학습관",
	},
	"111019": {
		District: "교육청",
		LibName:  "서울특별시교육청용산도서관",
	},
	"111020": {
		District: "교육청",
		LibName:  "서울특별시교육청정독도서관",
	},
	"111021": {
		District: "교육청",
		LibName:  "서울특별시교육청종로도서관",
	},
	"111022": {
		District: "교육청",
		LibName:  "서울특별시교육청노원평생학습관",
	},
	"111030": {
		District: "교육청",
		LibName:  "서울특별시교육청송파도서관",
	},
	"111031": {
		District: "교육청",
		LibName:  "서울특별시교육청마포평생학습관아현분관",
	},
	"111034": {
		District: "중랑구",
		LibName:  "중랑구립정보도서관",
	},
	"111035": {
		District: "성동구",
		LibName:  "성동구립도서관",
	},
	"111036": {
		District: "광진구",
		LibName:  "광진정보도서관",
	},
	"111037": {
		District: "강북구",
		LibName:  "강북문화정보도서관",
	},
	"111038": {
		District: "강북구",
		LibName:  "강북청소년문화정보도서관",
	},
	"111039": {
		District: "강남구",
		LibName:  "강남구립대치도서관",
	},
	"111040": {
		District: "금천구",
		LibName:  "금천구립독산도서관",
	},
	"111041": {
		District: "도봉구",
		LibName:  "도봉문화정보도서관",
	},
	"111042": {
		District: "은평구",
		LibName:  "은평구립도서관",
	},
	"111044": {
		District: "성북구",
		LibName:  "성북정보도서관",
	},
	"111045": {
		District: "관악구",
		LibName:  "관악중앙도서관",
	},
	"111047": {
		District: "노원구",
		LibName:  "노원어린이도서관",
	},
	"111048": {
		District: "성북구",
		LibName:  "아리랑도서관",
	},
	"111049": {
		District: "송파구",
		LibName:  "거마도서관",
	},
	"111051": {
		District: "서대문구",
		LibName:  "서대문구립이진아기념도서관",
	},
	"111052": {
		District: "강북구",
		LibName:  "솔샘문화정보도서관",
	},
	"111053": {
		District: "동대문",
		LibName:  "이문체육문화센터 어린이도서관",
	},
	"111056": {
		District: "강남구",
		LibName:  "강남구립즐거운도서관",
	},
	"111058": {
		District: "노원구",
		LibName:  "노원중앙도서관",
	},
	"111060": {
		District: "중랑구",
		LibName:  "중랑구립면목정보도서관",
	},
	"111062": {
		District: "관악구",
		LibName:  "글빛정보도서관",
	},
	"111063": {
		District: "관악구",
		LibName:  "성현동작은도서관",
	},
	"111065": {
		District: "동대문",
		LibName:  "동대문구정보화도서관",
	},
	"111067": {
		District: "강서구",
		LibName:  "강서구립길꽃어린이도서관",
	},
	"111068": {
		District: "성동구",
		LibName:  "성동구립 금호도서관",
	},
	"111069": {
		District: "강남구",
		LibName:  "강남구립청담도서관",
	},
	"111070": {
		District: "강남구",
		LibName:  "강남구립논현도서관",
	},
	"111071": {
		District: "강남구",
		LibName:  "논현문화마루도서관(분관)",
	},
	"111072": {
		District: "은평구",
		LibName:  "대조꿈나무어린이도서관",
	},
	"111073": {
		District: "강남구",
		LibName:  "강남구립정다운도서관",
	},
	"111075": {
		District: "강남구",
		LibName:  "강남구립행복한도서관",
	},
	"111076": {
		District: "성동구",
		LibName:  "성동구립 용답도서관",
	},
	"111077": {
		District: "금천구",
		LibName:  "금천구립가산도서관",
	},
	"111078": {
		District: "강동구",
		LibName:  "강동구립성내도서관",
	},
	"111081": {
		District: "양천구",
		LibName:  "신월음악도서관",
	},
	"111083": {
		District: "영등포구",
		LibName:  "대림도서관",
	},
	"111084": {
		District: "송파구",
		LibName:  "소나무언덕1호작은도서관",
	},
	"111085": {
		District: "노원구",
		LibName:  "월계도서관",
	},
	"111086": {
		District: "마포구",
		LibName:  "마포서강도서관",
	},
	"111097": {
		District: "강서구",
		LibName:  "강서구립푸른들청소년도서관",
	},
	"111098": {
		District: "강서구",
		LibName:  "강서구립꿈꾸는어린이도서관",
	},
	"111099": {
		District: "구로구",
		LibName:  "구로꿈나무어린이도서관",
	},
	"111100": {
		District: "구로구",
		LibName:  "꿈마을도서관",
	},
	"111101": {
		District: "강동구",
		LibName:  "강동구립해공도서관",
	},
	"111102": {
		District: "동작구",
		LibName:  "동작영어마루도서관",
	},
	"111103": {
		District: "강남구",
		LibName:  "강남구립역삼도서관",
	},
	"111104": {
		District: "마포구",
		LibName:  "꿈을이루는작은도서관",
	},
	"111105": {
		District: "마포구",
		LibName:  "늘푸른소나무작은도서관",
	},
	"111106": {
		District: "마포구",
		LibName:  "성메작은도서관",
	},
	"111107": {
		District: "동작구",
		LibName:  "약수도서관",
	},
	"111108": {
		District: "중구",
		LibName:  "가온도서관",
	},
	"111110": {
		District: "은평구",
		LibName:  "구립증산정보도서관",
	},
	"111111": {
		District: "영등포구",
		LibName:  "문래도서관",
	},
	"111113": {
		District: "금천구",
		LibName:  "금천구립금나래도서관",
	},
	"111114": {
		District: "송파구",
		LibName:  "소나무언덕2호작은도서관",
	},
	"111115": {
		District: "광진구",
		LibName:  "자양제4동도서관",
	},
	"111116": {
		District: "광진구",
		LibName:  "중곡문화체육센터도서관",
	},
	"111117": {
		District: "송파구",
		LibName:  "송파어린이도서관",
	},
	"111118": {
		District: "강서구",
		LibName:  "강서구립우장산숲속도서관",
	},
	"111119": {
		District: "성동구",
		LibName:  "성동구립 무지개도서관",
	},
	"111120": {
		District: "용산구",
		LibName:  "청파도서관",
	},
	"111121": {
		District: "영등포구",
		LibName:  "선유도서관",
	},
	"111123": {
		District: "성북구",
		LibName:  "해오름도서관",
	},
	"111124": {
		District: "은평구",
		LibName:  "구립응암정보도서관",
	},
	"111125": {
		District: "강동구",
		LibName:  "강동구립강일도서관",
	},
	"111126": {
		District: "은평구",
		LibName:  "구립상림도서관",
	},
	"111127": {
		District: "강남구",
		LibName:  "강남구립역삼푸른솔도서관",
	},
	"111128": {
		District: "구로구",
		LibName:  "온누리도서관",
	},
	"111129": {
		District: "도봉구",
		LibName:  "도봉구립학마을도서관",
	},
	"111131": {
		District: "노원구",
		LibName:  "파랑새 작은도서관",
	},
	"111132": {
		District: "은평구",
		LibName:  "은평작은도서관",
	},
	"111133": {
		District: "노원구",
		LibName:  "중계사랑 작은도서관",
	},
	"111135": {
		District: "노원구",
		LibName:  "수락 작은도서관",
	},
	"111136": {
		District: "서초구",
		LibName:  "반포1동 작은도서관",
	},
	"111137": {
		District: "은평구",
		LibName:  "신사어린이도서관",
	},
	"111138": {
		District: "서초구",
		LibName:  "반포본동 작은도서관",
	},
	"111140": {
		District: "서초구",
		LibName:  "잠원도서관",
	},
	"111141": {
		District: "서초구",
		LibName:  "방배본동 작은도서관",
	},
	"111142": {
		District: "서초구",
		LibName:  "방배2동 책사랑방",
	},
	"111146": {
		District: "서초구",
		LibName:  "방배4동 작은도서관",
	},
	"111148": {
		District: "동작구",
		LibName:  "동작샘터도서관",
	},
	"111160": {
		District: "관악구",
		LibName:  "은천동작은도서관",
	},
	"111163": {
		District: "서초구",
		LibName:  "방배1동 작은도서관",
	},
	"111171": {
		District: "마포구",
		LibName:  "아름드리작은도서관",
	},
	"111173": {
		District: "송파구",
		LibName:  "소나무언덕3호작은도서관",
	},
	"111174": {
		District: "강북구",
		LibName:  "송중문화정보도서관",
	},
	"111175": {
		District: "송파구",
		LibName:  "소나무언덕4호작은도서관",
	},
	"111176": {
		District: "강동구",
		LibName:  "강동구립암사도서관",
	},
	"111179": {
		District: "서대문구",
		LibName:  "남가좌새롬어린이도서관",
	},
	"111180": {
		District: "중구",
		LibName:  "남산타운 어린이도서관",
	},
	"111182": {
		District: "강북구",
		LibName:  "수유문화정보도서관",
	},
	"111188": {
		District: "성북구",
		LibName:  "종암동새날도서관",
	},
	"111189": {
		District: "도봉구",
		LibName:  "도봉아이나라도서관",
	},
	"111191": {
		District: "송파구",
		LibName:  "소나무언덕잠실본동도서관",
	},
	"111192": {
		District: "관악구",
		LibName:  "조원도서관",
	},
	"111208": {
		District: "도봉구",
		LibName:  "도봉구립무수골도서관",
	},
	"111215": {
		District: "구로구",
		LibName:  "하늘도서관",
	},
	"111216": {
		District: "구로구",
		LibName:  "개봉도서관",
	},
	"111217": {
		District: "관악구",
		LibName:  "관악산詩도서관",
	},
	"111218": {
		District: "관악구",
		LibName:  "낙성대공원도서관",
	},
	"111219": {
		District: "성북구",
		LibName:  "서경로꿈마루도서관",
	},
	"111220": {
		District: "중구",
		LibName:  "어울림도서관",
	},
	"111224": {
		District: "서초구",
		LibName:  "서초1동 작은도서관",
	},
	"111227": {
		District: "강남구",
		LibName:  "압구정동문고",
	},
	"111235": {
		District: "구로구",
		LibName:  "글마루한옥어린이도서관",
	},
	"111236": {
		District: "서초구",
		LibName:  "방배도서관",
	},
	"111239": {
		District: "노원구",
		LibName:  "화랑도서관",
	},
	"111240": {
		District: "양천구",
		LibName:  "개울건강도서관",
	},
	"111242": {
		District: "양천구",
		LibName:  "목마교육도서관",
	},
	"111243": {
		District: "서초구",
		LibName:  "반포4동 작은도서관",
	},
	"111244": {
		District: "서초구",
		LibName:  "양재2동 작은도서관",
	},
	"111246": {
		District: "서초구",
		LibName:  "양재1동 작은도서관",
	},
	"111247": {
		District: "서초구",
		LibName:  "반포3동 작은도서관",
	},
	"111249": {
		District: "서초구",
		LibName:  "반포2동 작은도서관",
	},
	"111250": {
		District: "서초구",
		LibName:  "서초3동 작은도서관",
	},
	"111251": {
		District: "서초구",
		LibName:  "서초4동 작은도서관",
	},
	"111252": {
		District: "서대문구",
		LibName:  "홍은도담도서관",
	},
	"111257": {
		District: "마포구",
		LibName:  "해오름작은도서관",
	},
	"111258": {
		District: "마포구",
		LibName:  "성산글마루작은도서관",
	},
	"111259": {
		District: "송파구",
		LibName:  "송파어린이영어작은도서관",
	},
	"111262": {
		District: "도봉구",
		LibName:  "방학1동 공립 작은도서관",
	},
	"111270": {
		District: "도봉구",
		LibName:  "도봉2동 공립 작은도서관",
	},
	"111271": {
		District: "도봉구",
		LibName:  "방학2동 공립 작은도서관",
	},
	"111272": {
		District: "도봉구",
		LibName:  "창1동 공립 작은도서관",
	},
	"111273": {
		District: "서대문구",
		LibName:  "아이누리작은도서관",
	},
	"111274": {
		District: "도봉구",
		LibName:  "쌍문3동 공립 작은도서관",
	},
	"111278": {
		District: "서대문구",
		LibName:  "행복작은도서관",
	},
	"111279": {
		District: "서대문구",
		LibName:  "문화촌작은도서관",
	},
	"111280": {
		District: "서대문구",
		LibName:  "꿈이있는작은도서관",
	},
	"111281": {
		District: "서대문구",
		LibName:  "알음알음작은도서관",
	},
	"111282": {
		District: "도봉구",
		LibName:  "쌍문2동 공립 작은도서관",
	},
	"111290": {
		District: "서대문구",
		LibName:  "새싹작은도서관",
	},
	"111299": {
		District: "은평구",
		LibName:  "불광천작은도서관",
	},
	"111301": {
		District: "성북구",
		LibName:  "석관동미리내도서관",
	},
	"111302": {
		District: "성북구",
		LibName:  "달빛마루도서관",
	},
	"111303": {
		District: "동대문",
		LibName:  "장안어린이도서관",
	},
	"111304": {
		District: "동대문",
		LibName:  "용두어린이영어도서관",
	},
	"111305": {
		District: "광진구",
		LibName:  "구의제3동도서관",
	},
	"111306": {
		District: "구로구",
		LibName:  "신도림어린이영어작은도서관",
	},
	"111307": {
		District: "노원구",
		LibName:  "상계도서관",
	},
	"111309": {
		District: "강서구",
		LibName:  "강서구립등빛도서관",
	},
	"111311": {
		District: "중랑구",
		LibName:  "중랑숲어린이도서관",
	},
	"111312": {
		District: "관악구",
		LibName:  "고맙습니다 하난곡작은도서관",
	},
	"111314": {
		District: "중구",
		LibName:  "서울도서관",
	},
	"111331": {
		District: "도봉구",
		LibName:  "창2동 공립 작은도서관",
	},
	"111344": {
		District: "성동구",
		LibName:  "성동구립 성수도서관",
	},
	"111346": {
		District: "마포구",
		LibName:  "마포어린이영어도서관",
	},
	"111347": {
		District: "강남구",
		LibName:  "강남구립도곡정보문화도서관",
	},
	"11134704": {
		District: "강남구",
		LibName:  "도곡2동작은도서관",
	},
	"11134706": {
		District: "강남구",
		LibName:  "개포4동작은도서관",
	},
	"11134707": {
		District: "강남구",
		LibName:  "일원본동작은도서관",
	},
	"11134711": {
		District: "강남구",
		LibName:  "신사동작은도서관",
	},
	"11134712": {
		District: "강남구",
		LibName:  "일원2동문고",
	},
	"11134713": {
		District: "강남구",
		LibName:  "수서동작은도서관",
	},
	"11134717": {
		District: "강남구",
		LibName:  "일원1동작은도서관",
	},
	"111348": {
		District: "도봉구",
		LibName:  "창4동 공립 작은도서관",
	},
	"111350": {
		District: "도봉구",
		LibName:  "창3동 공립 작은도서관",
	},
	"111351": {
		District: "도봉구",
		LibName:  "도봉1동 공립 작은도서관",
	},
	"111357": {
		District: "도봉구",
		LibName:  "방학3동 공립 작은도서관",
	},
	"111364": {
		District: "용산구",
		LibName:  "작은도서관 꿈꾸는책마을",
	},
	"111373": {
		District: "강서구",
		LibName:  "강서영어도서관",
	},
	"111374": {
		District: "마포구",
		LibName:  "복사골작은도서관",
	},
	"111377": {
		District: "서초구",
		LibName:  "서초구립반포도서관",
	},
	"111378": {
		District: "강서구",
		LibName:  "강서구립곰달래도서관",
	},
	"111380": {
		District: "성북구",
		LibName:  "정릉도서관",
	},
	"111383": {
		District: "강북구",
		LibName:  "미아문화정보도서관",
	},
	"111409": {
		District: "용산구",
		LibName:  "해다올 작은도서관",
	},
	"111414": {
		District: "노원구",
		LibName:  "메아리 작은도서관",
	},
	"111415": {
		District: "노원구",
		LibName:  "한울 작은도서관",
	},
	"111431": {
		District: "노원구",
		LibName:  "가재울지혜마루 작은도서관",
	},
	"111434": {
		District: "금천구",
		LibName:  "금천구립시흥도서관",
	},
	"111435": {
		District: "송파구",
		LibName:  "송파글마루도서관",
	},
	"111437": {
		District: "성북구",
		LibName:  "청수도서관",
	},
	"111439": {
		District: "동작구",
		LibName:  "사당솔밭도서관",
	},
	"111442": {
		District: "성동구",
		LibName:  "성동구립 청계도서관",
	},
	"111443": {
		District: "동작구",
		LibName:  "대방어린이도서관",
	},
	"111444": {
		District: "영등포구",
		LibName:  "여의샛강도서관",
	},
	"111445": {
		District: "동대문",
		LibName:  "동대문구답십리도서관",
	},
	"111446": {
		District: "동대문",
		LibName:  "휘경어린이도서관",
	},
	"111448": {
		District: "송파구",
		LibName:  "돌마리도서관",
	},
	"111452": {
		District: "은평구",
		LibName:  "구립은평뉴타운도서관",
	},
	"111453": {
		District: "강서구",
		LibName:  "강서구립가양도서관",
	},
	"111454": {
		District: "도봉구",
		LibName:  "둘리도서관",
	},
	"111456": {
		District: "송파구",
		LibName:  "가락몰도서관",
	},
	"111457": {
		District: "도봉구",
		LibName:  "도봉기적의도서관",
	},
	"111458": {
		District: "은평구",
		LibName:  "구립구산동도서관마을",
	},
	"111461": {
		District: "양천구",
		LibName:  "영어특성화도서관",
	},
	"111462": {
		District: "마포구",
		LibName:  "마포푸르메어린이도서관",
	},
	"111463": {
		District: "양천구",
		LibName:  "갈산도서관",
	},
	"111464": {
		District: "성북구",
		LibName:  "월곡꿈그림도서관",
	},
	"111465": {
		District: "강북구",
		LibName:  "삼각산어린이도서관",
	},
	"111466": {
		District: "서초구",
		LibName:  "서이도서관",
	},
	"111467": {
		District: "마포구",
		LibName:  "마포중앙도서관",
	},
	"111468": {
		District: "성북구",
		LibName:  "아리랑어린이도서관",
	},
	"111469": {
		District: "마포구",
		LibName:  "꿈나래어린이영어도서관",
	},
	"111470": {
		District: "강동구",
		LibName:  "강동구립천호도서관",
	},
	"111471": {
		District: "구로구",
		LibName:  "궁동어린이도서관",
	},
	"111472": {
		District: "송파구",
		LibName:  "송파위례도서관",
	},
	"111473": {
		District: "용산구",
		LibName:  "용산꿈나무도서관",
	},
	"111474": {
		District: "강남구",
		LibName:  "강남구립못골한옥어린이도서관",
	},
	"111476": {
		District: "서초구",
		LibName:  "서초구립내곡도서관",
	},
	"111477": {
		District: "서초구",
		LibName:  "서초그림책도서관",
	},
	"111478": {
		District: "중랑구",
		LibName:  "중랑구립양원숲속도서관",
	},
	"111479": {
		District: "강남구",
		LibName:  "강남구립못골도서관",
	},
	"111480": {
		District: "양천구",
		LibName:  "미감도서관",
	},
	"111481": {
		District: "양천구",
		LibName:  "해맞이역사도서관",
	},
	"111482": {
		District: "은평구",
		LibName:  "구립내를건너서숲으로도서관",
	},
	"111484": {
		District: "노원구",
		LibName:  "월계어린이도서관",
	},
	"111485": {
		District: "노원구",
		LibName:  "불암도서관",
	},
	"111487": {
		District: "성북구",
		LibName:  "성북길빛도서관",
	},
	"111488": {
		District: "중랑구",
		LibName:  "중랑구립상봉도서관",
	},
	"111489": {
		District: "구로구",
		LibName:  "구로기적의도서관",
	},
	"111490": {
		District: "동대문",
		LibName:  "배봉산 숲속도서관",
	},
	"111491": {
		District: "은평구",
		LibName:  "구립은뜨락도서관",
	},
	"111492": {
		District: "성북구",
		LibName:  "장위행복누림도서관",
	},
	"111493": {
		District: "서초구",
		LibName:  "서초구립양재도서관",
	},
	"111495": {
		District: "광진구",
		LibName:  "자양한강도서관",
	},
	"111496": {
		District: "도봉구",
		LibName:  "쌍문채움도서관",
	},
	"111498": {
		District: "중구",
		LibName:  "손기정 어린이도서관",
	},
	"111499": {
		District: "성동구",
		LibName:  "성동구립 매봉산 숲속도서관",
	},
	"111500": {
		District: "성북구",
		LibName:  "글빛도서관",
	},
	"111501": {
		District: "강동구",
		LibName:  "강동구립둔촌도서관",
	},
	"111502": {
		District: "동작구",
		LibName:  "김영삼도서관",
	},
	"111503": {
		District: "서초구",
		LibName:  "서초청소년도서관",
	},
	"111504": {
		District: "양천구",
		LibName:  "양천중앙도서관",
	},
	"111511": {
		District: "중구",
		LibName:  "신당누리도서관",
	},
	"111512": {
		District: "노원구",
		LibName:  "공릉행복도서관",
	},
	"111513": {
		District: "중구",
		LibName:  "다산성곽도서관",
	},
	"111514": {
		District: "마포구",
		LibName:  "마포소금나루도서관",
	},
	"111515": {
		District: "도봉구",
		LibName:  "김근태기념도서관",
	},
	"111516": {
		District: "동작구",
		LibName:  "까망돌도서관",
	},
	"111517": {
		District: "중구",
		LibName:  "손기정 문화도서관",
	},
	"111518": {
		District: "강남구",
		LibName:  "강남구립일원라온영어도서관",
	},
	"111519": {
		District: "노원구",
		LibName:  "하계어린이도서관",
	},
	"111521": {
		District: "광진구",
		LibName:  "군자동도서관",
	},
	"111522": {
		District: "동대문",
		LibName:  "휘경행복도서관",
	},
	"111523": {
		District: "광진구",
		LibName:  "아차산숲속도서관",
	},
	"111524": {
		District: "도봉구",
		LibName:  "원당마을한옥도서관",
	},
	"111526": {
		District: "강남구",
		LibName:  "강남구립개포하늘꿈도서관",
	},
	"141619": {
		District: "노원구",
		LibName:  "한내지혜의숲도서관",
	},
	"311034": {
		District: "노원구",
		LibName:  "노원구청 종합자료실",
	},
	"311744": {
		District: "관악구",
		LibName:  "용꿈꾸는작은도서관",
	},
	"511015": {
		District: "도봉구",
		LibName:  "북서울중학교",
	},
	"711002": {
		District: "동대문",
		LibName:  "이문어린이도서관",
	},
	"711023": {
		District: "용산구",
		LibName:  "효창동 작은도서관",
	},
	"711044": {
		District: "노원구",
		LibName:  "하늘 작은도서관",
	},
	"711073": {
		District: "서대문구",
		LibName:  "하늘샘작은도서관",
	},
	"711074": {
		District: "서대문구",
		LibName:  "북아현마을 북카페",
	},
	"711075": {
		District: "서대문구",
		LibName:  "햇살작은도서관",
	},
	"711076": {
		District: "서대문구",
		LibName:  "파랑새작은도서관",
	},
	"711077": {
		District: "서대문구",
		LibName:  "늘푸른열린작은도서관",
	},
	"711079": {
		District: "동대문",
		LibName:  "전곡마을 작은도서관",
	},
	"711080": {
		District: "동대문",
		LibName:  "장안 가온누리 작은도서관",
	},
	"711081": {
		District: "동대문",
		LibName:  "동대문구 장안 벚꽃길 작은도서관",
	},
	"711085": {
		District: "동대문",
		LibName:  "답십리2동 민들레 작은도서관",
	},
	"711086": {
		District: "동대문",
		LibName:  "이문1동 꿈꾸는 작은도서관",
	},
	"711088": {
		District: "동대문",
		LibName:  "휘경2동 꿈빛누리 작은도서관",
	},
	"711089": {
		District: "동대문",
		LibName:  "휘경1동 새싹마루 작은도서관",
	},
	"711090": {
		District: "중랑구",
		LibName:  "중화어린이도서관",
	},
	"711140": {
		District: "마포구",
		LibName:  "용강동작은도서관",
	},
	"711152": {
		District: "은평구",
		LibName:  "녹번만화도서관",
	},
	"711250": {
		District: "동대문",
		LibName:  "전농2동 뜨락 작은도서관",
	},
	"711251": {
		District: "동대문",
		LibName:  "장안1동 작은도서관",
	},
	"711265": {
		District: "강남구",
		LibName:  "강남구립열린도서관",
	},
	"711266": {
		District: "강남구",
		LibName:  "대치1동작은도서관",
	},
	"711267": {
		District: "강남구",
		LibName:  "삼성도서관",
	},
	"711268": {
		District: "강남구",
		LibName:  "세곡도서관",
	},
	"711269": {
		District: "강남구",
		LibName:  "역삼2동작은도서관",
	},
	"711302": {
		District: "은평구",
		LibName:  "은평어린이영어도서관",
	},
	"711304": {
		District: "은평구",
		LibName:  "응암1동문화의집문고",
	},
	"711305": {
		District: "은평구",
		LibName:  "효경골마을문고",
	},
	"711307": {
		District: "은평구",
		LibName:  "역촌누리작은도서관",
	},
	"711308": {
		District: "은평구",
		LibName:  "갈현1동문화의집",
	},
	"711310": {
		District: "은평구",
		LibName:  "은평구행정자료실",
	},
	"711311": {
		District: "동대문",
		LibName:  "답십리1동 아름드리 작은도서관",
	},
	"711316": {
		District: "도봉구",
		LibName:  "쌍문1동 공립 작은도서관",
	},
	"711317": {
		District: "도봉구",
		LibName:  "쌍문4동 공립 작은도서관",
	},
	"711318": {
		District: "도봉구",
		LibName:  "창5동 공립 작은도서관",
	},
	"711319": {
		District: "도봉구",
		LibName:  "도봉구청 행복작은도서관",
	},
	"711320": {
		District: "도봉구",
		LibName:  "지혜의등대 작은도서관",
	},
	"711323": {
		District: "도봉구",
		LibName:  "방학동육아종합지원센터",
	},
	"711324": {
		District: "도봉구",
		LibName:  "창동육아종합지원센터",
	},
	"711325": {
		District: "동작구",
		LibName:  "다울작은도서관",
	},
	"711352": {
		District: "동작구",
		LibName:  "국사봉숲속작은도서관",
	},
	"711359": {
		District: "양천구",
		LibName:  "방아다리문학도서관",
	},
	"711396": {
		District: "용산구",
		LibName:  "청파어린이영어도서관",
	},
	"711397": {
		District: "용산구",
		LibName:  "용암어린이영어도서관",
	},
	"711398": {
		District: "용산구",
		LibName:  "후암동 작은도서관 북앤캠프",
	},
	"711399": {
		District: "용산구",
		LibName:  "남영동 작은도서관",
	},
	"711400": {
		District: "용산구",
		LibName:  "원효로제2동 작은도서관",
	},
	"711401": {
		District: "용산구",
		LibName:  "한강로동 작은도서관",
	},
	"711402": {
		District: "용산구",
		LibName:  "이촌2동 작은도서관",
	},
	"711403": {
		District: "용산구",
		LibName:  "회나무 작은도서관",
	},
	"711404": {
		District: "용산구",
		LibName:  "한남동 작은도서관",
	},
	"711405": {
		District: "용산구",
		LibName:  "서빙고동 작은도서관",
	},
	"711410": {
		District: "동대문",
		LibName:  "이문2동 작은도서관",
	},
	"711427": {
		District: "동대문",
		LibName:  "동대문책마당도서관",
	},
	"711433": {
		District: "동대문",
		LibName:  "장안마루 작은도서관",
	},
	"711466": {
		District: "용산구",
		LibName:  "별밭 작은도서관",
	},
	"711468": {
		District: "서대문구",
		LibName:  "논골작은도서관",
	},
	"711485": {
		District: "동대문",
		LibName:  "장안꿈마루어린이 작은도서관",
	},
	"711492": {
		District: "강남구",
		LibName:  "세곡마루도서관",
	},
	"711497": {
		District: "용산구",
		LibName:  "청소년 푸르미르 작은도서관",
	},
	"711539": {
		District: "중구",
		LibName:  "장충동 작은도서관",
	},
	"711596": {
		District: "마포구",
		LibName:  "마포나루메타버스도서관",
	},
	"711603": {
		District: "노원구",
		LibName:  "꿈꾸는 작은도서관",
	},
	"711604": {
		District: "노원구",
		LibName:  "달내 작은도서관",
	},
	"711605": {
		District: "노원구",
		LibName:  "도란도란 작은도서관",
	},
	"711606": {
		District: "노원구",
		LibName:  "책이랑친구랑 작은도서관",
	},
	"711607": {
		District: "노원구",
		LibName:  "열린 작은도서관",
	},
	"711608": {
		District: "노원구",
		LibName:  "책누리 작은도서관",
	},
	"711609": {
		District: "노원구",
		LibName:  "책사랑 북카페",
	},
	"711610": {
		District: "노원구",
		LibName:  "푸른숲 작은도서관",
	},
	"711611": {
		District: "노원구",
		LibName:  "해솔 작은도서관",
	},
	"711612": {
		District: "노원구",
		LibName:  "반디 작은도서관",
	},
	"711613": {
		District: "노원구",
		LibName:  "가온 작은도서관",
	},
	"711614": {
		District: "노원구",
		LibName:  "노원문화원 작은도서관",
	},
	"711615": {
		District: "노원구",
		LibName:  "상계숲속 작은도서관",
	},
	"711616": {
		District: "노원구",
		LibName:  "초안산숲속 작은도서관",
	},
	"711617": {
		District: "노원구",
		LibName:  "한내행복발전소 작은도서관",
	},
	"711618": {
		District: "노원구",
		LibName:  "KB국민은행과 함께하는 나무 작은도서관",
	},
	"741562": {
		District: "노원구",
		LibName:  "향기나무도서관",
	},
}
