class Onyx < Formula
  desc "Cross-platform CLI for Onyx Cloud Database"
  homepage "https://github.com/OnyxDevTools/onyx-cli"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OnyxDevTools/onyx-cli/releases/download/v0.1.0/onyx_darwin_amd64.tar.gz"
      sha256 "197dceacdd3b7084a4f4604af4117ed94a3db9ebfa941cf53af2dbee7c799054"
    else
      url "https://github.com/OnyxDevTools/onyx-cli/releases/download/v0.1.0/onyx_darwin_arm64.tar.gz"
      sha256 "8e412e577493ce4226009fdf06bbde640d64e14b46e6620fd585ea4f2957929f"
    end
  end

  def install
    bin.install "onyx"
  end

  test do
    assert_match "onyx version", shell_output("#{bin}/onyx version")
  end
end
