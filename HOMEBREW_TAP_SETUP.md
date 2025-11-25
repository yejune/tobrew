# Homebrew Tap 설정 가이드

이 가이드는 `brew install yejune/tap/tobrew` 형태로 설치할 수 있게 Homebrew Tap을 설정하는 방법입니다.

## tobrew를 사용하면

tobrew는 이 모든 과정을 자동화합니다:

```bash
# 1. 프로젝트 설정 (한 번만)
tobrew init

# 2. 릴리스 (자동으로 버전 증가, 빌드, 태그, Formula 생성 및 배포)
tobrew release              # patch 버전 증가 (0.0.1 → 0.0.2)
tobrew release --minor      # minor 버전 증가 (0.0.2 → 0.1.0)
tobrew release --major      # major 버전 증가 (0.1.0 → 1.0.0)
```

**이 문서는 tobrew가 내부적으로 수행하는 작업을 설명합니다.** 수동으로 설정하려는 경우에만 아래 단계를 따르세요.

---

## 수동 설정 단계

### 1단계: homebrew-tap 저장소 생성

GitHub에 새 저장소를 만듭니다:
- 저장소 이름: `homebrew-tap` (반드시 `homebrew-` 접두사 필요)
- Public 저장소로 생성

```bash
# 로컬에서 저장소 생성
mkdir homebrew-tap
cd homebrew-tap
git init
echo "# Homebrew Tap" > README.md
git add README.md
git commit -m "Initial commit"
git remote add origin https://github.com/yejune/homebrew-tap.git
git branch -M main
git push -u origin main
```

### 2단계: 프로젝트 설정 파일 생성

프로젝트 루트에 `tobrew.yaml` 파일을 생성합니다:

```yaml
name: tobrew
description: "Automated Homebrew tap release tool for Go projects"
homepage: https://github.com/yejune/tobrew
license: MIT

github:
  user: yejune
  repo: tobrew
  tap_repo: homebrew-tap

build:
  command: go build -o build/{{.Name}} .

formula:
  install: |
    system "go", "build", "."
    bin.install "tobrew"

  test: |
    assert_match "tobrew version", shell_output("#{bin}/tobrew --version")

  caveats: |
    tobrew has been installed!

    Simple workflow:
      1. tobrew init              # Create config (once)
      2. tobrew release           # Release with patch bump
      3. tobrew release --minor   # Minor version bump
      4. tobrew release --major   # Major version bump

    Documentation: https://github.com/yejune/tobrew
```

### 3단계: 첫 릴리스 만들기

```bash
# 버전 태그 생성
git tag v0.0.1
git push origin v0.0.1
```

또는 GitHub UI에서:
1. Releases 탭 이동
2. "Create a new release" 클릭
3. Tag: `v0.0.1`
4. Title: `v0.0.1`
5. "Publish release" 클릭

### 4단계: SHA256 해시 계산

릴리스가 생성되면 tarball의 SHA256 해시를 계산:

```bash
# tarball 다운로드
curl -L -o tobrew-0.0.1.tar.gz \
  https://github.com/yejune/tobrew/archive/refs/tags/v0.0.1.tar.gz

# SHA256 계산 (macOS)
shasum -a 256 tobrew-0.0.1.tar.gz

# SHA256 계산 (Linux)
sha256sum tobrew-0.0.1.tar.gz
```

### 5단계: Homebrew Formula 생성

homebrew-tap 저장소에 `tobrew.rb` 파일 생성:

```ruby
class Tobrew < Formula
  desc "Automated Homebrew tap release tool for Go projects"
  homepage "https://github.com/yejune/tobrew"
  url "https://github.com/yejune/tobrew/archive/refs/tags/v0.0.1.tar.gz"
  sha256 "여기에_계산한_해시값_입력"
  license "MIT"
  head "https://github.com/yejune/tobrew.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "."
    bin.install "tobrew"
  end

  test do
    assert_match "tobrew version", shell_output("#{bin}/tobrew --version")
  end

  def caveats
    <<~EOS
      tobrew has been installed!

      Simple workflow:
        1. tobrew init              # Create config (once)
        2. tobrew release           # Release with patch bump
        3. tobrew release --minor   # Minor version bump
        4. tobrew release --major   # Major version bump

      Documentation: https://github.com/yejune/tobrew
    EOS
  end
end
```

### 6단계: Formula 커밋 및 푸시

```bash
cd homebrew-tap
git add tobrew.rb
git commit -m "Add tobrew formula v0.0.1"
git push
```

### 7단계: 설치 테스트

이제 사용자들이 다음 명령어로 설치할 수 있습니다:

```bash
# Tap 추가 (처음 한 번만)
brew tap yejune/tap

# 설치
brew install tobrew
```

또는 한 줄로:

```bash
brew install yejune/tap/tobrew
```

---

## tobrew를 사용한 자동화 예제

위 모든 과정은 tobrew로 자동화됩니다. 실제 tobrew 자체도 tobrew로 릴리스됩니다:

```bash
# tobrew 프로젝트에서
cd /path/to/tobrew

# 첫 설정 (tobrew.yaml 생성 - 이미 존재함)
tobrew init

# 새 버전 릴리스
tobrew release --minor

# 이 명령 하나가 자동으로:
# 1. tobrew.lock에서 현재 버전 읽기 (없으면 v0.0.0부터 시작)
# 2. 버전 증가 (v0.0.1 → v0.1.0)
# 3. 사용자 확인 (Continue? Y/n)
# 4. 프로젝트 빌드 (build.command 실행)
# 5. Git annotated 태그 생성 및 푸시
# 6. GitHub 처리 대기 (5초)
# 7. GitHub tarball 다운로드 및 SHA256 계산
# 8. Homebrew Formula 생성 (로컬에 .rb 파일 저장)
# 9. homebrew-tap 저장소 클론 (없으면 초기화)
# 10. Formula 파일 업데이트 및 커밋
# 11. git push (force 없이)
# 12. tobrew.lock 저장 (버전, 타임스탬프, SHA256)
```

## 업데이트 배포

### 수동 방식

새 버전을 릴리스할 때:

1. 프로젝트에서 새 태그 생성:
   ```bash
   git tag v0.0.2
   git push origin v0.0.2
   ```

2. SHA256 계산:
   ```bash
   curl -L -o tobrew-0.0.2.tar.gz \
     https://github.com/yejune/tobrew/archive/refs/tags/v0.0.2.tar.gz
   shasum -a 256 tobrew-0.0.2.tar.gz
   ```

3. homebrew-tap 저장소의 Formula 업데이트:
   - `url` 줄의 버전 업데이트 (v0.0.1 → v0.0.2)
   - 새 SHA256 계산 및 업데이트

4. 커밋 & 푸시:
   ```bash
   cd homebrew-tap
   git add tobrew.rb
   git commit -m "Update tobrew to v0.0.2"
   git push
   ```

### tobrew 방식

```bash
tobrew release  # 끝!
```

사용자들은 다음 명령어로 업데이트:
```bash
brew update
brew upgrade tobrew
```

## tobrew.lock 파일

tobrew는 `tobrew.lock` 파일로 버전을 추적합니다:

```yaml
version: v0.0.1
last_release: 2025-01-20T10:30:00Z
sha256: abc123...
```

이 파일은 자동으로 생성되고 업데이트되므로 수동으로 편집하지 마세요.

## 참고

- Homebrew Formula 문서: https://docs.brew.sh/Formula-Cookbook
- Tap 생성 가이드: https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap
- tobrew GitHub: https://github.com/yejune/tobrew
