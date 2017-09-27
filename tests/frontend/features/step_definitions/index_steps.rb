Given %r{^there are directories called "([^"]*)" and "([^"]*)"$} do |dir1, dir2|
  @directories = [dir1, dir2]
  @directories.each do |dir|
    expect(FileTest.exist?(File.join(REPO_PATH, dir)))
  end
end

Given %r{^they both contain Markdown files$} do
  @directories.each do |dir|
    expect(Dir.glob(File.join(REPO_PATH, dir, "*.md")).any?).to be true
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
    expect(page).to have_css(".card[data-filename='#{filename}']")
  end
end

Given %r{^I am on the "([^"]*)" index page$} do |name|
  path = File.join("/", "cms", name)
  visit(path)
  expect(page.current_path).to eql(path)
end

When %r{^I click the "([^"]*)" navigation link$} do |link_text|
  within("#application > .navbar") do
    page.find(".navbar-toggler-icon").click
    click_link(link_text)
  end
end

Then %r{^I should be on the "([^"]*)" index page$} do |name|
  expect(page.current_path).to eql("/cms/#{name}")
end

Then %r{^each directory index page should have the correct title:$} do |table|
	table.hashes.each do |row|
		visit("/cms/#{row['Directory']}")
		expect(page).to have_css("h2", text: row['Title'])
	end
end

Then %r{^I should be on the new document page for the '(.*?)' directory$} do |directory|
  expect(page).to have_css("h1", text: "New Document")
  expect(page.current_path).to eql("/cms/#{directory}/new")
end

Given %r{^there is no directory called "(.*?)"$} do |directory|
  expect(Dir.exist?(File.join(REPO_PATH, directory))).to be false
end