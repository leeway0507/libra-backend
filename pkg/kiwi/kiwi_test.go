package kiwi

import (
	"libra-backend/pkg/bm25"
	"log"
	"reflect"
	"slices"
	"testing"
)

func TestKiwi(t *testing.T) {
	kb := NewBuilder("./models/base", 1 /*=numThread*/, KIWI_BUILD_INTEGRATE_ALLOMORPH /*=options*/)

	k := kb.Build()
	defer k.Close() // don't forget to Close()!
	t.Run("version", func(t *testing.T) {
		Version()
	})

	t.Run("Analyze", func(t *testing.T) {
		keyword := "도커, 컨테이너 빌드업! Docker, container build-up! : 최적의 컨테이너 서비스를 위한 도커 활용법"
		func() {
			tokens, err := k.Analyze(keyword, 1 /*=topN*/, KIWI_MATCH_NORMALIZE_CODA)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("tokens: %#+v\n", tokens)
		}()

	})
	t.Run("Analyze_Noun", func(t *testing.T) {
		keyword := "도커, 컨테이너 빌드업! Docker, container build-up! : 최적의 컨테이너 서비스를 위한 도커 활용법"
		func() {
			tokens, err := k.Analyze_Noun(keyword, 1 /*=topN*/, KIWI_MATCH_NORMALIZE_CODA)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("tokens: %#+v\n", tokens)
		}()
	})
	t.Run("Remove Duplicates", func(t *testing.T) {
		tokens := []string{"도커", "컨테이너", "빌드업", "최적", "컨테이너", "서비스", "도커", "활용법"}
		candidate := k.RemoveDuplicates(tokens)
		result := []string{"서비스", "활용법", "도커", "컨테이너", "빌드업", "최적"}

		slices.Sort(candidate)
		slices.Sort(result)
		if !reflect.DeepEqual(candidate, result) {
			log.Printf("candidate: %#+v\n", candidate)
			log.Printf("result: %#+v\n", result)
			log.Fatal("candidate doesn't match with result")
		}

	})

	t.Run("BM25", func(t *testing.T) {
		corpus := []string{
			"파이썬 프로그래밍",
			"시간순삭 파이썬",
			"파이썬을 이용한 알고리즘 트레이딩 =  : 아이디어 도출부터 클라우드 배포까지",
			"Hello IT 파이썬을 제대로 활용해보려고 해",
			"그림으로 배우는 파이썬 기초 문법",
			"쓸모 있는 파이썬 프로그램 40개",
			"(Better Python Code) 파이썬 코딩의 기술 51",
			"(파이썬과 함께하는)생활 속의 통계학:초보 통계학 탐험가를 위한 입문서",
			"파이썬을 활용한 소프트웨어 아키텍처:견고하고 확장 가능한 애플리케이션 아키텍처 설계",
			"파이썬 프로그래밍:데이터 분석 프로젝트로 프로그래밍 사고력 기르기",
			"(Step by step) 파이토치 딥러닝 프로그래밍 .PART 1",
			"(Do it!) 점프 투 파이썬 =  : 이미 200만 명이 '점프 투 파이썬'으로 프로그래밍을 시작했다!",
			"딱 한 줄로! 파이썬 제대로 코딩하기 ",
			"코스페이시스 with 파이썬",
			"비전공자를 위한 이해할 수 있는 파이썬:AI 시대에 최적화된 파이썬 공부법",
			"파이토치 트랜스포머를 활용한 자연어 처리와 컴퓨터비전 심층학습 ",
			"파이썬으로 배우는 포트폴리오  = Portfolio with python   ",
			"(코딩 자율학습)나도코딩의 파이썬 입문",
			"알고리즘 인사이드 with 파이썬 =  : 86개 풀이로 문제 해결 능력, 사고력을 키우는 알고리즘 & 자료구조 입문서",
			"(IT 비전공자를 위한)파이썬 업무 자동화(RPA)",
			"제대로 빠르게 파이썬 입문 ",
			"파이썬 딥러닝으로 시작하는 이상 징후 탐지(한국어판)",
			"프로그래머를 위한 수학 (파이썬으로 하는 3D 그래픽스, 시뮬레이션, 머신러닝)",
			"(파이썬을 이용한) 비트코인 자동매매 :실시간 자동매매 시스템 개발 입문 ",
			"파이썬 생활밀착형 프로젝트 =  : 웹 크롤링, 카카오톡 메시지 보내기, 업무자동화까지 11가지 파이썬 프로젝트",
			"(Do it!)점프 투 파이썬",
			"파이썬 GUI 프로그래밍 쿡북 3/e ",
			"(Do it!)점프 투 파이썬",
			"개발자를 위한 실전 선형대수학",
			"(생성형 AI 빅3)챗GPT·미드저니·스테이블 디퓨전 =  : 알수록 진짜 돈 되는 기술",
			"시작하세요! Final Cut Pro 10.6 - 빠르크의 3분 강좌와 함께하는 파이널 컷 프로 유튜브 영상 제작",
			"파이썬 크래시 코스",
			"(데이터 과학을 위한) 파이썬과 R:오픈소스를 활용한 데이터 분석, 시각화, 머신러닝",
			"다빈치 리졸브 :아직도 돈내고 영상 편집하니? ",
			"파이썬 데브옵스 프로그래밍 :파이썬으로 하는 인프라 자동화 ",
			"금융 데이터 분석을 위한 파이썬 판다스",
			"파이썬 클린 코드 =  : 유지보수가 쉬운 파이썬 코드를 만드는 비결",
			"사장님 몰래하는 파이썬 업무 자동화",
			"(다양한 캐글 예제와 함께 기초 알고리즘부터 최신 기법까지 배우는)파이썬 머신러닝 완벽 가이드",
			"파이썬 웹 프로그래밍",
			"파이썬과 비교하며 배우는 러스트 프로그래밍 가장 사랑받는 언어 러스트를 배우는 가장 확실한 방법",
			"(파이썬과 비교하며 배우는)러스트 프로그래밍 : 가장 사랑받는 언어 러스트를 배우는 가장 확실한 방법",
			"Hey, 파이썬! 생성형 AI 활용 앱 만들어 줘",
			"파이썬 텍스트 마이닝 바이블",
			"파이썬으로 하는 마케팅 연구와 분석 : 데이터 처리부터 시각화까지",
			"(파이썬을 이용한)퀀트 투자 포트폴리오 만들기 =  ",
			"데이터 경영을 위한 파이썬 =  : 성공하는 CEO의 시스템 분석 툴",
			"(Do it!)점프 투 장고 = Do it! jump to Django  : 파이썬 웹 개발부터 배포까지!",
			"Do it! 쉽게 배우는 파이썬 데이터 분석 : 데이터 분석 프로젝트 전 과정 수록!",
			"데이터 분석으로 배우는 파이썬 문제 해결 : 부동산 데이터 분석부터 AWS 아키텍처 구축, 대시보드 제작까지",
			"100만 뷰 프리미어 프로들의 유튜브 영상편집 테크닉",
		}

		kb := NewBuilder("./models/base", 1 /*=numThread*/, KIWI_BUILD_INTEGRATE_ALLOMORPH /*=options*/)

		k := kb.Build()
		defer k.Close() // don't forget to Close()!

		tokenizer := func(s string) []string {
			tokens, err := k.Analyze(s, 1 /*=topN*/, KIWI_MATCH_ALL)
			if err != nil {
				log.Println("tokenizer error", err)
			}
			return tokens
		}

		bm25, err := bm25.NewBM25Okapi(corpus, tokenizer, 0.5, 0.1, nil)
		if err != nil {
			t.Fatal(err)
		}
		query := "점프 파이썬"
		kb.AddWord(query, "NNP", 0)
		tokenizedQuery := tokenizer(query)
		// for _, v := range corpus {
		// 	log.Printf("%#+v\n", tokenizer(v))
		// }

		x, err := bm25.GetScores(tokenizedQuery)
		if err != nil {
			t.Fatal(err)
		}

		for i, v := range x {
			log.Printf("score: %v title: %s", v, corpus[i])
		}

	})

}
