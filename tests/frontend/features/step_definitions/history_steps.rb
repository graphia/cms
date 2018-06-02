Given %r{^I am on the appendix history page for "([^"]*)"$} do |arg1|
  path = "/cms/appendices/appendix_1/history"
  visit(path)
  expect(page.current_path).to eql(path)
end

When %r{^I navigate to the history page$} do
  path = "/cms/history"
  visit(path)
  expect(page.current_path).to eql(path)
end

When %r{^I am on the history page$} do
  path = "/cms/history"
  visit(path)
  expect(page.current_path).to eql(path)
end


Given %r{^I am on the document history page for "([^"]*)"$} do |arg1|
  path = "/cms/documents/document_1/history"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^there is a document with '(\d+)' revisions$} do |revs|

  g = Git.open(REPO_PATH)

    # Set some committer info
    g.config('user.name', 'Roy Snyder')
    g.config('user.email', 'roy.snyder@springfield.court.us')

    new_document = <<~CONTENTS
      ---
      title: History Test
      author: Joey Crusher
      ---

      Revision 1
    CONTENTS

    @doc = "history"

    dir = File.join(REPO_PATH, "documents", @doc)
    Dir.mkdir(dir)
    full_path = File.join(dir, "index.md")
    File.write(full_path, new_document)

    g.add(all: true)

    g.commit("Added revision 1")

    # get the hash of the latest commit (the one right above!)
    @commits =  [{name: "revision 1", hash: g.log.first.to_s}]

    2.upto(revs.to_i) do |i|
      File.write(full_path, File.read(full_path).gsub(/Revision \d+/, "Revision #{i}"))
      g.add(all: true)
      g.commit("Added revision #{i}")
      @commits.push({name: "revision #{i}", hash: g.log.first.to_s})
    end

    expect(@commits.length).to eql(revs.to_i)

end

When %r{^I am on the document's history page$} do
  path = "/cms/documents/#{@doc}/history"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see each revision listed$} do
  within("ol.commit-list") do
    @commits.each do |commit|
      expect(page).to have_css(".card.history.commit-#{commit[:hash]}")
    end
  end
end

Given %r{^there is a document with some revisions$} do
  step "there is a document with '2' revisions"
end

When %r{^I click the '(.*?)' button for a particular revision$} do |button_text|
  @commit = @commits[1]
  within(".card.history.commit-#{@commit[:hash]}") do
    page.find(".btn", text: button_text).click
  end
end

Then %r{^I should be on that commit's page$} do
  expect(page.current_path).to eql("/cms/commits/#{@commit[:hash]}")
end

Then %r{^that revision's diff should be visible beneath the revision entry$} do
  within(".card.history.commit-#{@commit[:hash]}") do
    expect(page).to have_css("pre.diff")
  end
end

Then %r{^the diff should have correctly marked insertions and deletions$} do
  within(".card.history.commit-#{@commit[:hash]}") do
    expect(page).to have_css("pre.diff del", text: 1)
    expect(page).to have_css("pre.diff ins", text: 2)
  end
end


Then %r{^I should see my commit$} do
  expect(page).to have_css("div[data-commit-hash='#{@hash}']")
end

Given %r{^I have added a new file with a multiline commit message$} do

  @summary = "Do the Bartman"
  @message_1 = "Yo! Hey, what's happening, dude?"
  @message_2 = "I'm the guy with the rep for being rude."

  @commit_message = <<~MSG
    #{@summary}

    #{@message_1}
    #{@message_2}
  MSG

  step "I have added a new file"
end

Then %r{^the first line of my commit message should be a header$} do
  within("div[data-commit-hash='#{@hash}']") do
    expect(page).to have_css("h4", text: @summary)
  end
end

Then %r{^all subsequent lines should be paragraphs$} do
  within("div[data-commit-hash='#{@hash}']") do
    [@message_1, @message_2].each do |m|
      expect(page).to have_css("p", text: m)
    end
  end
end

Then %r{^the commit should include the committer name$} do
  within("div[data-commit-hash='#{@hash}']") do
    expect(page).to have_css(".author", text: @commit_author)
  end
end