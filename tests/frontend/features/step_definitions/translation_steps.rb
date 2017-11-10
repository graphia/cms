Given %r{^my document has been translated into 'Finnish' and 'Swedish'$} do
  write_english_file
  write_finnish_file
  write_swedish_file
  @document = "translated_doc"
end

Given %r{^my document has been translated into 'Finnish'$} do
  write_english_file
  write_finnish_file
  @document = "translated_doc"
end

Given %r{^my document has been translated into 'Swedish'$} do
  write_english_file
  write_swedish_file
  @document = "translated_doc"
end

When %r{^I visit the documents index page$} do
  path = "/cms/documents"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see my document listed with 'Finnish' and 'Swedish' flags$} do
  within("div[data-filename='translated_doc'] .card-footer") do
    expect(page).to have_css("li[data-lang='Finnish']", text: "🇫🇮")
    expect(page).to have_css("li[data-lang='Swedish']", text: "🇸🇪")
  end
end

Given %r{^my document has not been translated into any other languages$} do
  @document = "document_1"
end

Then %r{^my document shouldn't have any alternative languages section$} do
  within(".card[data-filename='#{@document}']") do
    expect(page).not_to have_css(".card-footer")
  end
end

# Used on the document index feature
When %r{^I click the 'Finnish' link$} do
  within(".card[data-filename='#{@document}'] li[data-lang='Finnish']") do
    link = page.find("a")
    scroll_into_view(link)
    link.click
  end
end

# Used on the document index feature
Then %r{^I should be on the document's '(.*?)' translation$} do |lang|
  codes = {
    "Finnish" => "fi",
    "Swedish" => "sv"
  }
  expect(page.current_path).to eql("/cms/documents/#{@document}.#{codes[lang]}.md")
end

# Used on the document summary feature
When %r{^I click the 'Swedish' link$} do

  within("a[data-filename='#{@document}'] li[data-lang='Swedish']") do
    link = page.find("a")
    scroll_into_view(link)
    link.click
  end
end

Given %r{^I navigate to my document's 'show' page$} do
  @document = "document_1.md"
  path = "/cms/documents/#{@document}"
  visit(path)
  expect(page.current_path).to eql(path)
end


Given %r{^my document is untranslated$} do
  # do nothing
end

When %r{^I click the 'Translate' dropdown button$} do
  within(".document-metadata") do
    click_button("Translate")
  end
end

Then %r{^I should see a list of available languages:$} do |table|
  within(".document-metadata .translation-options") do
    table.transpose.raw.flatten.each do |language|
      expect(page).to have_css("button", text: language)
    end
  end
end

Then %r{^the existing language '(.*?)' should not be included$} do |language|
  within(".document-metadata .translation-options") do
    expect(page).not_to have_css("button", text: language)
  end
end

Then %r{^the existing languages '(.*?)' and '(.*?)' should not be included$} do |language1, language2|
  within(".document-metadata .translation-options") do
    [language1, language2].each do |language|
        expect(page).not_to have_css("button", text: language)
    end
  end
end

Given %r{^there is already a 'Swedish' translation of my document$} do
  write_swedish_file("document_1")
end

When %r{^I click the '(.*?)' translation option$} do |language|
  within(".document-metadata .translation-options") do
    click_button(language)
  end
end

Then %r{^I should be on the new 'Swedish' document$} do
  expect(page).to have_css(".alert.alert-success")
  expect(page.current_path).to eql("/cms/documents/document_1.sv.md")
end

Then %r{^I should see my document listed under '(.*?)'$} do |dir|
  within("div[data-directory='#{dir.downcase}']") do
    expect(page).to have_content("A noble spirit embiggens the smallest man.")
  end
end

Then %r{^it should have 'Finnish' and 'Swedish' flags$} do
  within("a[data-filename='#{@document}']") do
    expect(page).to have_css("li[data-lang='Finnish']", text: "🇫🇮")
    expect(page).to have_css("li[data-lang='Swedish']", text: "🇸🇪")
  end
end

def write_translated_file(name, email, contents, code, message, filename="translated_doc")
  Git.open(REPO_PATH).tap do |g|
    g.config('user.name', name)
    g.config('user.email', email)
    File.write(
      File.join(
        REPO_PATH,
        "documents",
        [filename, code, "md"].compact.join(".")
      ),
      contents
    )
    g.add(all: true)
    g.commit(message)
  end
end

def write_english_file(filename="translated_doc")
  @english_title = "A noble spirit embiggens the smallest man."
  contents = <<~CONTENTS
    ---
    title: #{@english_title}
    ---

    We paid a short guy to write it,
    But we never saw him again.
    The tune we stole from the French.
    There's a few things they do well.
  CONTENTS

  write_translated_file(
    'Melvin van Horne',
    'mvh@hbowtime.com',
    contents,
    nil,
    "English file added",
    filename
  )
end

def write_finnish_file(filename="translated_doc")
  contents = <<~CONTENTS
    ---
    title: Finland's Culture (Model U.N. Club)
    ---

    Oi maamme, Suomi, synnyinmaa,
    soi, sana kultainen!
  CONTENTS

  write_translated_file(
    'B. A. Bäckström',
    'ba.bs@moomin.fi',
    contents,
    "fi",
    "Finnish file added",
    filename
  )
end

def write_swedish_file(filename="translated_doc")
  contents = <<~CONTENTS
    ---
    title: Du gamla, Du fria, Du fjällhöga nord
    ---

    Du gamla, Du fria, Du fjällhöga nord
    Du tysta, Du glädjerika sköna!
  CONTENTS

  write_translated_file(
    'Sven Simpson',
    'sven.simpson@ikea.se',
    contents,
    "sv",
    "Swedish file added",
    filename
  )
end