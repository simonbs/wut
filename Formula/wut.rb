class Wut < Formula
  desc "Git worktree manager that keeps worktrees out of your repo"
  homepage "https://github.com/simonbs/wut"
  url "https://github.com/simonbs/wut/archive/refs/tags/v0.1.1.tar.gz"
  sha256 "582f06836f28dfe7a621ea24cac4d277709c139abf5a50717839c369f16d267a"
  head "https://github.com/simonbs/wut.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/wut"
  end

  test do
    assert_match "wut", shell_output("#{bin}/wut --help")
  end
end
