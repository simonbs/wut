class Wut < Formula
  desc "Git worktree manager that keeps worktrees out of your repo"
  homepage "https://github.com/simonbs/wut"
  url "https://github.com/simonbs/wut/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "7f01da8bfa34bb390cf5c2ed8d962df3f958cf78b3b8bef91fbb23926185e469"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/wut"
  end

  def caveats
    <<~EOS
      🚀 To enable shell integration, add to your ~/.zshrc or ~/.bashrc:

        eval "$(wut init)"

      Then restart your shell or run: source ~/.zshrc
    EOS
  end

  test do
    assert_match "wut", shell_output("#{bin}/wut --help")
  end
end
