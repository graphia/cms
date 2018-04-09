Given %r{^I am on the homepage$} do
  path = "/cms/"
  visit(path)
  expect(page.current_path).to eql(path)
end

When %r{^I visit the homepage$} do
  step "I am on the homepage"
end

Then %r{^I should be on the homepage$} do
  expect(page.current_path).to eql("/cms")
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
    doc = File.dirname(file).split("/").last
    FileUtils.cp(file, File.join(REPO_PATH, "documents", filename))
    g.add(all: true)
    g.commit_all("Added #{doc}")
  end

end

Then %r{^the recent changes summary should contain a list of commits$} do

  within(".recent-updates") do
    ["Added a", "Added b", "Added c"].each do |message|
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
    expect(page).to have_css("h4 > a", text: name)
  end
end

Given %r{^the documents directory contains the following files:$} do |table|
  table.transpose.raw.flatten.each do |doc|
    expect(File.exist?(File.join(REPO_PATH, "documents", doc, "index.md"))).to be true
  end
end

Then %r{^I should see all three documents listed$} do
  within(".documents") do
    ["document 1", "document 2", "document 3"].each do |title|
      expect(page).to have_css("a[data-filename='#{title.gsub(' ', '_')}']", text: title)
    end
  end
end

Then %r{^there should be a '(.*?)' button$} do |text|
  within(".documents") do
    expect(page).to have_css("a.btn", text: text)
  end
end

Given %r{^the '(.*?)' directory contains no files$} do |dir|
  # This could be improved with Dir.children, but need a newer Ruby
  # than 2.4.0 https://ruby-doc.org/core-2.4.1/Dir.html
  expect(Dir.entries(File.join(REPO_PATH, "empty"))).to eql(['.', '..', '_index.md'])
end

Then %r{^I see a 'no files' alert in the empty section$} do
  within(".empty") do
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

    These *documents* are **amazing**, aren't they?
    FM
  )
end

Then %r{^I should see the custom description$} do
  within(".documents") do
    expect(page).to have_css(".directory-description", text: "Documents go here")
  end
end

Then %r{^I should see the custom title$} do
  expect(page).to have_css("h4", text: "Important Documents")
end

Then %r{^I should see a "(.*?)" section$} do |name|
  expect(page).to have_css(".#{name.downcase}", text: name.capitalize)
end

Given %r{^there is one user$} do
  # Do nothing, there's only one user set up by default
end

Then %r{^the statistics panel's "(.*?)" count should equal "(.*?)"$} do |stat, count|
  within(".card.statistics") do
    expect(page).to have_css(".#{stat}-count", text: count)
  end
end

Then %r{^I should see "(.*?)" files of type "(.*?)"$} do |count, label|
  within(".card.statistics table.file-statistics > tbody") do
    expect(page).to have_css("tr[data-count-type='#{label.downcase}'] > td", text: label)
    expect(page).to have_css("tr[data-count-type='#{label.downcase}'] > td", text: count)
  end
end

Given %r{^my company name is "(.*?)"$} do |company_name|

  visit "/cms"
  token = evaluate_script("localStorage.token")

  uri = URI('http://127.0.0.1:9095/api/server_info')
  req = Net::HTTP::Get.new(uri)
  req['Authorization'] = "Bearer #{token}"

  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end

  json = JSON.parse(res.body)

  @config_site_title = json["title"]
  expect(@config_site_title).to eql(company_name)

end

Then %r{^the header in the main navigation should contain "(.*?)"$} do |company_name|
  expect(page).to have_css("nav .navbar-brand", text: company_name)
end

Then %r{^the title should be a link to the homepage$} do
  within("nav.navbar") do
    expect(page).to have_link(@config_site_title, href: "/cms")
  end
end