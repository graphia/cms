Given %r{^a repository has been initialised$} do
  setup_repo
end

Given %r{^there is an empty directory in place of a repository$} do
  FileUtils.rm_rf(REPO_PATH)
  FileUtils.mkdir(REPO_PATH)
end