def modify_file(path, filename)
  @modified_document = <<~CONTENTS
    ---
    title: Bart's Friend Falls in Love
    author: Jay Kogen & Wallace Wolodarsky
    ---
    # We started out like Kirk and Luann…

    …but it ended in tragedy.
  CONTENTS

  Dir.mkdir(path) unless Dir.exist?(path)
  File.write(File.join(path, filename), @modified_document)
end


Given %r{^I have added a new file$} do
  g = Git.open(REPO_PATH)

  @commit_author = "Roy Snyder"

  # Set some committer info
  g.config('user.name', @commit_author)
  g.config('user.email', 'roy.snyder@springfield.court.us')

  @new_document = <<~CONTENTS
    ---
    title: Bart's Friend Falls in Love
    author: Jay Kogen & Wallace Wolodarsky
    ---
    # We started out like Romeo and Juliet…

    …but it ended in tragedy.
  CONTENTS

  dir = File.join(REPO_PATH, "documents", "s03e23")
  Dir.mkdir(dir)
  File.write(File.join(dir, "index.md"), @new_document)

  g.add(all: true)

  @commit_message ||= "Added Bart's Friend Falls in Love"
  g.commit(@commit_message)

  # get the hash of the latest commit (the one right above!)
  @hash = g.log.first.to_s

end

When %r{^I navigate to the commit's details page$} do
  path = "/cms/commits/#{@hash}"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see some information specific to the commit$} do
  expect(page).to have_css("dd", text: @hash)
  expect(page).to have_css("dd > a", text: "Roy Snyder")
  expect(page).to have_css("dd", text: "Added Bart's Friend Falls in Love")
end

Then %r{^I should see a '(.*)' section with the file's path for a title$} do |context|
  contexts = {
    "blue" => ["file-updated", "border-info"],
    "green" => ["file-created", "border-success"],
    "red" => ["file-deleted", "border-danger"]
  }
  expect(page).to have_css("div.card.#{contexts[context].join(".")}")
end

Then %r{^it should contain the entire file's contents$} do
  within("div.card") do
    expect(page.find("pre")).to have_content(@new_document)
  end
end

Given %r{^I have made changes to an existing file$} do
  step "I have added a new file"

  g = Git.open(REPO_PATH)

  # Set some committer info
  g.config('user.name', 'Constance Harm')
  g.config('user.email', 'constance.harm@springfield.court.us')

  modify_file(File.join(REPO_PATH, "documents", "s03e23"), "index.md")

  g.add(all: true)
  g.commit("Switched Romeo and Juliet for Kirk and Luann")

  # get the hash of the latest commit (the one right above!)
  @hash = g.log.first.to_s
end

Given %r{^I have deleted a file$} do
  step "I have added a new file"

  g = Git.open(REPO_PATH)

  # Set some committer info
  g.config('user.name', 'Lionel Hutz')
  g.config('user.email', 'lionel.hutz@lawyers101.com')

  g.remove(File.join("documents", "s03e23", "index.md"))
  g.commit("Deleted Bart's Friend Falls in Love ")

  @hash = g.log.first.to_s
end

Then %r{^it should contain a colourised diff showing changes made$} do
  within("div.card.file-updated") do
    expect(page.find("pre")).to have_content("We started out like")
    expect(page.find("pre ins")).to have_content("Kirk and Luann")
    expect(page.find("pre del")).to have_content("Romeo and Juliet")
  end
end

Then %r{^the diff '(.*)' icon should be visible$} do |context|
  within("div.card.commit-file h2") do
    expect(page).to have_css("svg.octicon-diff-#{context}")
  end
end

Given %r{^I have modified one file and removed another in a single commit$} do

  step "I have added a new file"
  g = Git.open(REPO_PATH)

  # Set some committer info
  g.config('user.name', 'Lenny Leonard')
  g.config('user.email', 'lenny.leonard@nuclear.springfield.com')

  modify_file(File.join(REPO_PATH, "documents", "s03e23"), "index.md")

  g.remove(File.join("documents", "document_1", "index.md"))

  g.add(all: true)
  g.commit("Various changes")
  @hash = g.log.first.to_s
end

Then %r{^I should see two file sections, one for each affected file$} do
  within("ol.files") do
    expect(page).to have_css("div.card.commit-file", count: 2)
  end
end

Then %r{^the commits's page title should contain "(.*?)" and the short hash$} do |title|
  expect(page).to have_title("#{title} #{@hash.slice(0..6)}")
end

Given %r{^I have added, deleted and modified images in one commit$} do

  g = Git.open(REPO_PATH)

  # Add a file
  FileUtils.cp(
    File.join(IMAGES_PATH, "bart.png"),
    File.join(REPO_PATH, "documents", "document_1", "images", "bart.png")
  )

  # Delete a file
  FileUtils.rm(File.join(REPO_PATH, "appendices", "appendix_1", "images", "image_1.png"))

  # Modify a file
  FileUtils.cp(
    File.join(IMAGES_PATH, "lisa.jpg"),
    File.join(REPO_PATH, "documents", "document_2", "images", "image_1.jpg")
  )

  g.add(all: true)
  g.commit("Various changes")
  @hash = g.log.first.to_s

end

Then %r{^I should see the added file displayed in a green box$} do
  within(".card.file-created") do
    img = page.find("div.bg-success > img")
    expect(img['src']).to start_with("data:image/png;base64,")
  end
end

Then %r{^I should see the removed file displayed in a red box$} do
  within(".card.file-deleted") do
    img = page.find("div.bg-danger > img")
    expect(img['src']).to start_with("data:image/png;base64,")
  end
end

Then %r{^I should see the modified file in a green box next to the old version in a red box$} do
  within(".card.file-updated") do

    of = page.find("div.bg-danger > img")
    expect(of['src']).to start_with("data:image/png;base64,")

    nf = page.find("div.bg-success > img")
    expect(nf['src']).to start_with("data:image/png;base64,")

  end
end

Then %r{^the "(.*?)" file should have tooltip "(.*?)"$} do |operation, text|
  operations = {
    "added" => "file-created",
    "updated" => "file-updated",
    "deleted" => "file-deleted"
  }

  within("div.card.#{operations[operation]}") do
    # the bootstrap tooltip plugin moves the attr from title to data-original-title
    expect(page).to have_css("h2[data-toggle='tooltip'][data-original-title='#{text}']")
  end
end