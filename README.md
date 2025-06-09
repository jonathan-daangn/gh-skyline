# gh-skyline

GitHub 활동 그래프를 기반으로 3D STL 파일을 생성하는 CLI 도구입니다. 커스텀 텍스트, 기간, 월별 범위 등 다양한 옵션을 지원합니다.

---

## 설치 및 실행

Go 1.23 이상이 필요합니다.

```bash
git clone <이 레포지토리 주소>
cd gh-skyline
go mod tidy
```

### 실행 예시

```bash
go run main.go --user jjogeon --top-text "@jjogeon" --right-text "Special Thanks 2025" --year "2022" --start-month 3 --end-month 5
```

---

## 주요 옵션

- `--user`         : 기여자 GitHub 아이디 (기본값: 인증된 사용자)
- `--year`         : 연도 또는 연도 범위 (예: 2022, 2019-2022)
- `--full`         : 가입 연도부터 현재까지 전체 그래프 생성
- `--top-text`     : STL 상단에 표시할 텍스트
- `--right-text`   : STL 우측에 표시할 텍스트
- `--start-month`  : 시작 월 (기본값: 1)
- `--end-month`    : 종료 월 (기본값: 12)
- `--output`       : 출력 파일명 지정 (기본값: 자동 생성)
- `--art-only`     : STL 파일 없이 ASCII 아트만 출력
- `--web`          : 해당 사용자의 GitHub 프로필을 브라우저로 열기
- `--debug`        : 디버그 로그 출력

---

## 사용 예시

- 기본 사용 (현재 인증된 사용자, 올해 기준):
  ```bash
  go run main.go
  ```
- 특정 사용자/연도:
  ```bash
  go run main.go --user jjogeon --year 2022
  ```
- 커스텀 텍스트/월 범위:
  ```bash
  go run main.go --user jjogeon --top-text "@jjogeon" --right-text "Special Thanks 2025" --year "2022" --start-month 3 --end-month 5
  ```
- 가입~현재까지 전체 그래프:
  ```bash
  go run main.go --user jjogeon --full
  ```
- ASCII 아트만 출력:
  ```bash
  go run main.go --art-only
  ```

---

## 여러 명의 STL을 한 번에 생성하는 예시 (Wintertech Internship 2025)

아래와 같이 여러 명의 STL을 각각 생성할 수 있습니다:

```bash
go run main.go --top-text "Hoji" --user "wlgh1553" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5    
go run main.go --top-text "Kiru" --user "kiru211" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Jared" --user "bbang3" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Tavian" --user "ocahs9" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Wade.kim" --user "wadekim2880" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Jenna.park" --user "gahyuun" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Linker" --user "Yoon-Hae-Min" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Vinci" --user "2hyunbin" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Jodie" --user "devyubin" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Bobae" --user "cobinding" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Uno" --user "hnnynh" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Aina" --user "Aina-an" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Gilroy" --user "gil-roy" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Erik" --user "Erik-Kim" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Kevin.kim" --user "khyojun" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Sydney.lee" --user "rheeri" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
go run main.go --top-text "Maya" --user "ChaeAg" --right-text "Wintertech Internship 2025" --start-month 3 --end-month 5
```

---

## 빌드

```
