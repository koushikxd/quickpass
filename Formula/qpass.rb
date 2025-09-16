class Onepass < Formula
  desc "Secure password generator CLI tool"
  homepage "https://github.com/koushik/qpass"
  url "https://github.com/koushik/qpass/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "YOUR_SHA256_HASH_HERE"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    assert_match "Usage of", shell_output("#{bin}/qpass --help")
    
    system "#{bin}/qpass", "-n", "8"
  end
end
