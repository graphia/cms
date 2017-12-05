Given %r{^a repository has been initialised$} do
  setup_repo
end

Given %r{^there is an empty directory in place of a repository$} do
  FileUtils.rm_rf(REPO_PATH)
  FileUtils.mkdir(REPO_PATH)
  expect{Git.open(REPO_PATH)}.to raise_error(ArgumentError, "path does not exist")
end

Then %r{^a repository should have been initialised$} do
  expect(Git.open(REPO_PATH)).to be_a(Git::Base)
end