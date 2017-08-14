Given %r{^the document I'm working on already has an attachment$} do

  @image_path = File.join(REPO_PATH, "documents", "document_2", "images", "image_1.jpg")

  # explicitly call the RSpec matcher here as it clashes with Capybara's `all`
  expect([
    FileTest.exist?(File.join(REPO_PATH, "documents", "document_2.md")),
    FileTest.exist?(@image_path)
  ]).to RSpec::Matchers::BuiltIn::All.new(be true)
end

Then %r{^I should see the image in the gallery$} do

  image = File.open(@image_path)

  uri = "data:image/jpeg;base64,#{Base64.strict_encode64(image.read)}"

  within(".gallery") do
    expect(page).to have_css("img")
    expect(uri).to eql(page.find("img")[:src])
  end
end

Given %r{^I am on the edit page for my document$} do
  steps %{
    Given I am on the edit document page for "document_2.md"
  }
end

When %r{^I drag an image from the gallery to the editor$} do
  pending # Write code here that turns the phrase above into concrete actions
end

Then %r{^the markdown image placeholder should be added to the editor$} do
  pending # Write code here that turns the phrase above into concrete actions
end