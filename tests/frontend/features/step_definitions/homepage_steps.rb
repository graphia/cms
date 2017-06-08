Given %r{^I am on the homepage$} do
  path = "/cms/"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see a summary of recent changes$} do
  expect(page).to have_css("h4", text: "Recent Updates")
end

Then %r{^I should see a statistics section$} do
  expect(page).to have_css("h4", text: "Statistics")
end

Given %r{^there have been some recent commits$} do

  g = Git.open(REPO_PATH)

  # Set some committer info
  g.config('user.name', 'Dewey Largo')
  g.config('user.email', 'dewey.largo@springfield.k12.us')

  # Add some files
  Dir.glob(File.join(SAMPLES_PATH, "**/*.md")) do |file|
    filename = File.basename(file)
    FileUtils.cp(file, File.join(REPO_PATH, "documents", filename))
    g.add(all: true)
    g.commit_all("Added #{filename}")
  end

end

Then %r{^the recent changes summary should contain a list of commits$} do
  within(".recent-updates") do
    ["Added a.md", "Added b.md", "Added c.md"].each do |message|
      expect(page).to have_content(message)
    end

    expect(page).to have_content(
      "Dewey Largo",
      count: Dir.glob(File.join(SAMPLES_PATH, "**/*.md")).count
    )
  end
end