class Wut < Formula
  desc "Git worktree manager that keeps worktrees out of your repo"
  homepage "https://github.com/simonbs/wut"
  url "https://github.com/simonbs/wut/archive/refs/tags/v0.3.1.tar.gz"
  sha256 "e2a99626f2a47d6fa588643ed5fa6ba7691d7764af9a36a92da7c18d710d6813"
  head "https://github.com/simonbs/wut.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/wut"
  end

  test do
    assert_match "wut", shell_output("#{bin}/wut --help")
  end
end
