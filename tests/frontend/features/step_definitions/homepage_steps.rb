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

Given %r{^the following directories exist in the repository$} do |table|
  table.transpose.raw.flatten.each do |name|
    expect(Dir.exist?(File.join(REPO_PATH, name))).to be true
  end
end

Then %r{^I should see a section for each directory$} do
  ["Appendices", "Important Documents"].each do |name|
    expect(page).to have_css("h4.card-header", text: name)
  end
end

Given %r{^the documents directory contains the following files:$} do |table|

    %w{document_1.md document_2.md document_3.md}.each do |name|
    expect(File.exist?(File.join(REPO_PATH, "documents", name))).to be true
  end
end

Then %r{^I should see all three documents listed$} do
  within(".card.documents") do
    {
      "document_1.md" => "document 1",
      "document_2.md" => "document 2",
      "document_3.md" => "document 3"
    }.each do |filename, title|
      expect(page).to have_css("a[data-filename='#{filename}']", text: title)
    end
  end
end

Then %r{^there should be a 'new file' button$} do
  within(".card.documents") do
    expect(page).to have_css("a.btn", text: "Create a document")
  end
end

Given %r{^the '(.*?)' directory contains no files$} do |dir|
  # This could be improved with Dir.children, but need a newer Ruby
  # than 2.4.0 https://ruby-doc.org/core-2.4.1/Dir.html
  expect(Dir.entries(File.join(REPO_PATH, "empty"))).to eql(['.', '..', '_index.md'])
end

Then %r{^I see a 'no files' alert in the operating procedures section$} do
  within(".card.empty") do
    expect(page).to have_css("div.alert", text: "There's nothing here yet")
  end
end

Given %r{^the 'documents' directory has title and description metadata$} do
  contents = File.read(File.join(REPO_PATH, "documents", "_index.md"))

  expect(contents).to eql(
    <<~FM
    ---
    title: Important Documents
    description: Documents go here
    ---
    FM
  )
end

Then %r{^I should see the custom description$} do
  within(".card.documents") do
    expect(page).to have_css(".card-body", text: "Documents go here")
  end
end

Then %r{^I should see the custom title$} do
  expect(page).to have_css("h4.card-header", text: "Important Documents")
end