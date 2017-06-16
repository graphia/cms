Given %r{^I have added a new file$} do
  g = Git.open(REPO_PATH)

  # Set some committer info
  g.config('user.name', 'Roy Snyder')
  g.config('user.email', 'roy.snyder@springfield.court.us')

  @new_document = <<~CONTENTS
    ---
    title: Bart's Friend Falls in Love
    author: Jay Kogen & Wallace Wolodarsky
    ---
    # We started out like Romeo and Juliet…

    …but it ended in tragedy.
  CONTENTS

  File.write(File.join(REPO_PATH, "documents", "s03e23.md"), @new_document)

  g.add(all: true)


  g.commit("Added Bart's Friend Falls in Love")

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
    "blue" => ["file-updated", "card-outline-info"],
    "green" => ["file-created", "card-outline-success"],
    "red" => ["file-deleted", "card-outline-danger"]
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

  @modified_document = <<~CONTENTS
    ---
    title: Bart's Friend Falls in Love
    author: Jay Kogen & Wallace Wolodarsky
    ---
    # We started out like Kirk and Luann…

    …but it ended in tragedy.
  CONTENTS

  File.write(File.join(REPO_PATH, "documents", "s03e23.md"), @modified_document)

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

  g.remove(File.join("documents", "s03e23.md"))
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