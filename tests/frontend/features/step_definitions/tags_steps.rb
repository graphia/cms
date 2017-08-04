Given %r{^I have selected the tags editor$} do
  within(".frontmatter-fields") do
    @tags_editor = page.find(".tags-input input")
  end
end

When %r{^I enter the following tags and press "([^"]*)":$} do |key, table|
  table.transpose.raw.flatten.each do |tag|
    @tags_editor.native.send_keys(tag.chars, key.to_sym)
  end
end

Then %r{^I should see the following tags listed:$} do |table|
  within(".frontmatter-fields") do
    table.transpose.raw.flatten.each do |tag|
      expect(page).to have_css("span.tag", text: tag)
    end
  end
end

Then %r{^I add tags for Sales and Marketing$} do
  steps %{
    Given I have selected the tags editor
    When I enter the following tags and press "enter":
      | Sales     |
      | Marketing |
  }
end

Then %r{^I should see my document with the correct tags$} do
  within(".document-meta") do
    %w{Sales Marketing}.each do |tag|
      expect(page).to have_css("span.tag", text: tag)
    end
  end
end