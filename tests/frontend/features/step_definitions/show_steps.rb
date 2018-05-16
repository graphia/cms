Given %r{^a document called '(.*?)' exists$} do |name|
  expect(File.exist?(File.join(REPO_PATH, "appendices", name, "index.md"))).to be true
end


Given %r{^I am on the show page for a non\-existant document$} do
  path = "/cms/appendices/does-not-exist"
  visit(path)
  expect(page.current_path).to eql(path)
end

When %r{^I navigate to that document's 'show' page$} do
  path = "/cms/appendices/appendix_1"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I am on the show page for a document with metadata$} do
  path = "/cms/documents/document_1"
  visit(path)
  expect(page.current_path).to eql(path)
end

And %r{^I am on the document's show page$} do
  step "I navigate to that document's 'show' page"
end

Then %r{^I should see the correctly\-formatted document$} do
  expect(page).to have_css("h1", text: "Appendix 1")
  expect(page).to have_css("p", text: "Lorem ipsum dolor sit amet")
end

When %r{^the document has some frontmatter set up$} do
  contents = File.readlines(File.join(REPO_PATH, "appendices", "appendix_1", "index.md"))

  fm_range = contents
    .each
    .with_index
    .with_object([]) do |(line, i), bounds|
      bounds << i if line == "---\n"
    end

  fm = contents[fm_range[0]+1 .. fm_range[1]-1].join

  @yaml = YAML.load(fm).to_h

end

Then %r{^I should see the following frontmatter items with correct values:$} do |table|
  table.transpose.raw.flatten.each do |section|
    expect(@yaml[section.downcase].to_s).not_to be_empty

    within(".card.document-metadata") do
      if @yaml[section.downcase].is_a?(Array)
        @yaml[section.downcase].each do |list_item|
          expect(page).to have_css(".tag.badge", text: list_item)
        end
      else
        expect(page).to have_css("dd", text: @yaml[section.downcase])
      end
    end
  end
end

Then %r{^I should see a toolbar with the following buttons:$} do |table|
  table.transpose.raw.flatten do |button_text|
    within(".document-metadata .btn-toolbar") do
      expect(page).to have_css(".btn", text: button_text)
    end
  end
end

When %r{^I click the toolbar's '(.*?)' button$} do |button_text|
  within(".document-metadata .btn-toolbar") do
    page.find('.btn', text: button_text).click
  end
end

Then %r{^I should be on the document's edit page$} do
  expect(page).to have_css('h1', text: "Appendix 1")
  expect(page.current_path).to eql("/cms/appendices/appendix_1/edit")
end

Then %r{^I should be on the directory's index page$} do
  expect(page).to have_css('h1', text: 'Appendices')
  expect(page.current_path).to eql("/cms/appendices")
end

Then %r{^the document should have been deleted$} do
  expect(File.exist?(File.join(REPO_PATH, "appendices", "appendix_1.md"))).to be false
end

Then %r{^I should be on the document's history page$} do
  expect(page.current_path).to eql("/cms/appendices/appendix_1/history")
end

When %r{^I visit my document's 'English' page$} do
  path = "/cms/documents/translated_doc"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^my document should have links to 'English', 'Finnish' and 'Swedish' in the Translations section$} do

  langs = {
    "en" => "🇬🇧",
    "sv" => "🇸🇪",
    "fi" => "🇫🇮"
  }
  within(".translations") do
    langs.each do |lang, flag|
      expect(page).to have_css("li[data-lang='#{lang}']", text: flag)
    end
  end

end