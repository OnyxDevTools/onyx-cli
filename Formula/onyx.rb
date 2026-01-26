class Onyx < Formula
  desc "Cross-platform CLI for Onyx Cloud Database"
  homepage "https://github.com/OnyxDevTools/onyx-cli"
  version "0.2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/OnyxDevTools/onyx-cli/releases/download/v0.1.0/onyx_darwin_amd64.tar.gz"
      sha256 "6eefd28836715203f50d2d5f22d1bab7ea19ae260f684fbec633365a0751e2ce"
    else
      url "https://github.com/OnyxDevTools/onyx-cli/releases/download/v0.1.0/onyx_darwin_arm64.tar.gz"
      sha256 "7e686ec0b983d34caffbd56e96f35479c85fe0088436144c1064be3489d64909"
    end
  end

  def install
    bin.install "onyx"
  end

  test do
    assert_match "onyx version", shell_output("#{bin}/onyx version")
  end
end
