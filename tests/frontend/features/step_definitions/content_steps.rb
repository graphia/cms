Then %r{^I should see a button with text '(.*?)'$} do |caption|
  expect(page).to have_css("button.btn-primary", text: caption)
end

Then %r{^I should see text '(.*?)'$} do |text|
  expect(page).to have_content(text)
end

Then %r{^there should be an alert with the message "(.*)"$} do |message|
  page.save_screenshot("/tmp/out.png")
  expect(page).to have_css(".alert", text: message, wait: 5)
end

Then %r{^the main heading should be "(.*?)"$} do |heading|
  expect(page.first("h1")).to have_content(heading)
end

Then %r{^the navigation bar should contain the following links:$} do |table|
  within("nav.navbar") do
    table.transpose.raw.flatten.each do |link|
      expect(page).to have_css("li > a", text: link)
    end
  end
end

Then %r{^the page's title (?:should be|is) "(.*?)"$} do |title|
  expect(page).to have_title(title)
end

Then %r{^now the title is "(.*?)"$} do |title|
  expect(page).to have_title(title)
end