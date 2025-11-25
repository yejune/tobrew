# tobrew

**tobrew** - Go 프로젝트를 위한 자동화된 Homebrew tap 릴리스 도구

한 번의 명령으로 Homebrew tap 릴리스를 자동화하세요. 더 이상 수동 버전 관리, SHA256 계산, tap 저장소 업데이트가 필요 없습니다.

## 주요 기능

- ✅ **자동 버전 관리** - `tobrew.lock`이 현재 버전 추적
- 🚀 **원클릭 릴리스** - `tobrew release` 하나로 모든 것 해결
- 📝 **다양한 설정 포맷** - YAML, JSON, TOML 지원
- 🔐 **자동 SHA256 계산** - GitHub 릴리스에서 자동 계산
- 🍺 **Homebrew formula 생성** - 자동으로 formula 파일 생성
- 📦 **자동 tap 업데이트** - homebrew-tap 저장소 자동 관리
- 🎯 **간단한 워크플로우** - 단 2개 명령어만 필요

## 설치

### 방법 1: go install 사용 (권장)

```bash
go install github.com/yejune/tobrew@latest
```

### 방법 2: 소스에서 빌드

```bash
git clone https://github.com/yejune/tobrew.git
cd tobrew
go build
./tobrew install
```

### 업데이트

```bash
tobrew self-update
```

## 빠른 시작

### 1. 설정 초기화 (한 번만)

Go 프로젝트 디렉토리에서:

```bash
tobrew init
```

`tobrew.yaml` 파일이 생성됩니다. JSON이나 TOML도 사용 가능:

```bash
tobrew init --format json
tobrew init --format toml
```

### 2. 설정 편집

`tobrew.yaml`을 편집하고 정보를 업데이트하세요:

```yaml
name: myapp
description: "나의 멋진 CLI 도구"
homepage: https://github.com/username/myapp
license: MIT

github:
  user: username           # GitHub 사용자명
  repo: myapp             # 저장소 이름
  tap_repo: homebrew-tap  # tap 저장소 이름 (반드시 "homebrew-"로 시작)

build:
  command: go build -o build/{{.Name}} .

formula:
  install: |
    system "go", "build", "."
    bin.install "myapp"

  test: |
    assert_match "myapp", shell_output("#{bin}/myapp --version")

  caveats: |
    myapp가 설치되었습니다!
    'myapp --help'를 실행하여 시작하세요.
```

### 3. GitHub tap 저장소 생성

GitHub에서 `homebrew-tap`이라는 새 저장소를 생성하세요 (반드시 `homebrew-`로 시작).

### 4. 릴리스!

```bash
# 첫 릴리스 (v0.0.1 생성)
tobrew release

# Patch 릴리스 (v0.0.1 → v0.0.2)
tobrew release

# Minor 릴리스 (v0.0.2 → v0.1.0)
tobrew release --minor

# Major 릴리스 (v0.1.0 → v1.0.0)
tobrew release --major
```

다음 작업이 자동으로 수행됩니다:
1. **로드**: `tobrew.lock`에서 현재 버전 읽기
2. **증가**: 버전 증가 (patch/minor/major)
3. **빌드**: 프로젝트 빌드
4. **태그**: git tag 생성 및 푸시
5. **다운로드**: 릴리스 tarball 다운로드 및 SHA256 계산
6. **생성**: Homebrew formula 생성
7. **업데이트**: homebrew-tap 저장소 업데이트
8. **저장**: 새 버전을 `tobrew.lock`에 저장

이제 사용자들은 다음과 같이 설치할 수 있습니다:

```bash
brew install username/tap/myapp
```

## 명령어

### `tobrew init`

새 설정 파일을 초기화합니다.

```bash
tobrew init                    # tobrew.yaml 생성
tobrew init --format json      # tobrew.json 생성
tobrew init --format toml      # tobrew.toml 생성
tobrew init -o custom.yaml     # 커스텀 출력 경로
```

### `tobrew release`

자동 버전 증가와 함께 릴리스를 생성합니다.

```bash
tobrew release              # Patch: v1.0.0 → v1.0.1 (기본값)
tobrew release --patch      # Patch: v1.0.0 → v1.0.1 (명시적)
tobrew release --minor      # Minor: v1.0.1 → v1.1.0
tobrew release --major      # Major: v1.1.0 → v2.0.0
```

### `tobrew sync`

lock 파일을 원격 git 태그와 동기화합니다.

```bash
tobrew sync
```

다음과 같은 경우에 유용합니다:
- lock 파일이 실제 릴리스와 동기화되지 않은 경우
- 다른 머신에서 작업하는 경우
- 실패한 릴리스에서 복구하는 경우

## 버전 관리

tobrew는 `tobrew.lock` 파일을 사용하여 프로젝트 버전을 추적합니다:

```yaml
version: v1.2.3
last_release: 2025-11-25T15:30:00+09:00
sha256: abc123...
fingerprint: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```

- **첫 릴리스**: `v0.0.1`에서 시작
- **자동 증가**: 버전 번호를 직접 지정할 필요 없음
- **시맨틱 버저닝**: semver를 따름 (MAJOR.MINOR.PATCH)
- **Git 추적**: `tobrew.lock`을 저장소에 커밋
- **자동 동기화**: fingerprint가 다르거나 (다른 머신) 태그 충돌 시 자동으로 원격과 동기화

## 일반적인 워크플로우

```bash
# 초기 설정 (한 번만)
cd myproject
tobrew init
# tobrew.yaml 편집
git add tobrew.yaml tobrew.lock
git commit -m "Add tobrew config"

# 개발 사이클
# ... 코드 변경 ...
git add -A
git commit -m "Add new feature"
git push

# 릴리스!
tobrew release              # Patch 릴리스
# 또는
tobrew release --minor      # Minor 릴리스
# 또는
tobrew release --major      # Major 릴리스
```

## 실전 예제: tobrew 자체

이 프로젝트는 자기 자신을 릴리스하기 위해 tobrew를 사용합니다!

[tobrew.yaml.example](tobrew.yaml.example) 파일을 참고하세요.

## 왜 tobrew를 사용해야 하나요?

[goreleaser](https://goreleaser.com/) 같은 다른 도구들은 포괄적이지만 복잡합니다. tobrew는:

- ✅ **더 간단함** - 단 2개 명령어 (`init`과 `release`)
- ✅ **집중적** - Homebrew tap 관리에만 집중
- ✅ **자동적** - 수동 입력 없이 버전 관리
- ✅ **가벼움** - 최소한의 설정만 필요

Homebrew 배포만 필요한 Go CLI 도구에 완벽합니다.

## 문제 해결

### "failed to download tarball"

- GitHub에 git tag가 존재하는지 확인
- 태그를 푸시한 후 몇 초 대기
- GitHub 저장소가 public이거나 접근 권한이 있는지 확인

### "tap update failed"

- `homebrew-tap` 저장소가 존재하는지 확인
- tap 저장소에 푸시 권한이 있는지 확인
- 저장소 이름이 `homebrew-`로 시작하는지 확인

### "invalid version format"

- `tobrew.lock` 파일의 버전이 올바른지 확인 (예: `v1.2.3`)
- `tobrew.lock`을 삭제하면 `v0.0.1`부터 새로 시작

## 라이선스

MIT License - 자세한 내용은 LICENSE 파일 참조

## 기여

기여를 환영합니다! 이슈나 PR을 열어주세요.

## 관련 프로젝트

- [docker-bootapp](https://github.com/yejune/docker-bootapp) - tobrew를 사용하는 실전 예제
- [goreleaser](https://goreleaser.com/) - 포괄적인 릴리스 자동화 도구

---

더 쉬운 Homebrew 릴리스를 위해 ❤️를 담아 만들었습니다
