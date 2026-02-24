class Wut < Formula
  desc "Git worktree manager that keeps worktrees out of your repo"
  homepage "https://github.com/simonbs/wut"
  url "https://github.com/simonbs/wut/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "3c12bb96666c5bc69e94bb2a59d3cabb912e13c23f49e4164fdff253931cacb8"
  head "https://github.com/simonbs/wut.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/wut"
  end

  test do
    assert_match "wut", shell_output("#{bin}/wut --help")
  end
end
