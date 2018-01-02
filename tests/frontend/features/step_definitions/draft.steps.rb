Then %r{^the "(.*?)" checkbox should be checked$} do |name|
  within("label", text: /^#{name}$/) do
    expect(page.find("input[type='checkbox']")).to be_checked
  end
end

Then %r{^the "(.*?)" checkbox should be unchecked$} do |name|
  within("label", text: /^#{name}$/) do
    expect(page.find("input[type='checkbox']")).not_to be_checked
  end
end

When %r{^I fill in the rest of the document form and submit it$} do
  steps %{
    And I enter some text into the editor
    And I fill in the document metadata
    And I submit the form
  }
end

Then %r{^my document should( not)? be a draft$} do |negate|

  draft_content = negate ? "draft: false" : "draft: true"

  # wait upto a second for the file to be present
  1.upto(10) do
    begin
      file = File.read(File.join(REPO_PATH, "documents", "sample-document.md"))

      expect(file).not_to be_empty
      expect(file).to have_content(draft_content)
      break

    rescue Errno::ENOENT
      sleep(0.1)
    end
  end

end