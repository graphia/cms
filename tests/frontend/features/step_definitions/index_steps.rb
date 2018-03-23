Given %r{^there are directories called "([^"]*)" and "([^"]*)"$} do |dir1, dir2|
  @directories = [dir1, dir2]
  @directories.each do |dir|
    expect(FileTest.exist?(File.join(REPO_PATH, dir)))
  end
end

Given %r{^they both contain Markdown files$} do
  @directories.each do |dir|
    expect(Dir.glob(File.join(REPO_PATH, dir, "**/index.md")).any?).to be true
  end
end

When %r{^I navigate to the "([^"]*)" index page$} do |name|
  path = File.join("/", "cms", name)
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see a list containing the contents of the "(.*?)" directory$} do |dir|
  files = Dir
    .glob(File.join(REPO_PATH, dir, "*.md"))
    .map{|path| File.basename(path)}
    .reject {|filename| filename == "_index.md"}

  files.each do |filename|
    filename_without_ext = filename.gsub(".md", "")
    expect(page).to have_css(".card[data-filename='#{filename_without_ext}']")
  end
end

Given %r{^I am on the "([^"]*)" index page$} do |name|
  path = File.join("/", "cms", name)
  visit(path)
  expect(page.current_path).to eql(path)
end

When %r{^I click the "([^"]*)" navigation link$} do |link_text|
  within("header > .navbar") do

    # If the window is too small and the nav is collapsed,
    # we need to toggle the nav first. Normally, this is
    # unnecessary, but here's how to do it:
    #
    #page.find(".navbar-toggler-icon").click

    click_link(link_text)
  end
end

Then %r{^I should be on the "([^"]*)" index page$} do |name|
  expect(page.current_path).to eql("/cms/#{name}")
end

Then %r{^each directory index page should have the correct title:$} do |table|
	table.hashes.each do |row|
		visit("/cms/#{row['Directory']}")
		expect(page).to have_css("h1", text: row['Title'])
	end
end

Then %r{^I should be on the new document page for the '(.*?)' directory$} do |directory|
  expect(page).to have_css("h1", text: "New Document")
  expect(page.current_path).to eql("/cms/#{directory}/new")
end

Given %r{^there is no directory called "(.*?)"$} do |directory|
  expect(Dir.exist?(File.join(REPO_PATH, directory))).to be false
end

Given %r{^I have some documents that are drafts$} do
  file = File.read(File.join(REPO_PATH, "appendices", "appendix_1", "index.md"))
  expect(file).to have_content("draft: true")
end

Then %r{^the draft document should be highlighted$} do
  within(".document-list") do
    # highlighted by a border-warning with data-draft attr and matching filename stub
    matcher = "div.border-warning[data-filename='appendix_1'][data-draft='true']"
    expect(page).to have_css(matcher)
  end
end

When %r{^I click the link to "(.*)"$} do |name|
  within(".document-list") { click_link name }
end

Then %r{^I should be on the "([^"]*)" show page$} do |arg1|
  expect(page.current_path).to eql("/cms/documents/document_1")
end