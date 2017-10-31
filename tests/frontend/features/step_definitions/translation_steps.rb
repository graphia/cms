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

When %r{^I click the 'Finnish' link$} do
  within(".card[data-filename='#{@document}'] li[data-lang='Finnish']") do
    link = page.find("a")
    scroll_into_view(link)
    link.click
  end
end

Then %r{^I should be on the document's 'Finnish' translation$} do
  expect(page.current_path).to eql("/cms/documents/#{@document}.fi.md")
end

def write_translated_file(name, email, contents, code, message)
  Git.open(REPO_PATH).tap do |g|
    g.config('user.name', name)
    g.config('user.email', email)
    File.write(
      File.join(
        REPO_PATH,
        "documents",
        ["translated_doc", code, "md"].compact.join(".")
      ),
      contents
    )
    g.add(all: true)
    g.commit(message)
  end
end

def write_english_file
  contents = <<~CONTENTS
    ---
    title: A noble spirit embiggens the smallest man.
    ---

    We paid a short guy to write it,
    But we never saw him again.
    The tune we stole from the French.
    There's a few things they do well.
  CONTENTS

  write_translated_file(
    'B. A. Bäckström',
    'ba.bs@moomin.fi',
    contents,
    nil,
    "Finnish file added"
  )
end

def write_finnish_file
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
    "Finnish file added"
  )
end

def write_swedish_file
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
    "Swedish file added"
  )
end