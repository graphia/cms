Given %r{^I am on my ssh key settings page$} do
  path = "/cms/settings/keys"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I have no SSH keys$} do
  # do nothing
end

Given %r{^I have an SSH key$} do
  steps %{
		Given I am on my ssh key settings page
		And I enter 'home-desktop' in the 'Name' field
		And I paste my public SSH key into the 'ssh-key' text area
		When I submit the form
		Then there should be an alert with the message "SSH Key Created"
  }
end

When %r{^I visit my ssh key settings page$} do
  step "I am on my ssh key settings page"
end

Then %r{^I should see an alert informing me why I might want to add one$} do
  within(".no-keys") do
    expect(page).to have_css("h4.alert-heading", text: "You have no keys")
    expect(page).to have_content("In order to copy the entire repository to your own machine")
  end
end

Given %r{^I paste my public SSH key into the 'ssh\-key' text area$} do
  key = File.read("../backend/certificates/valid.pub").strip
  # fill_in "ssh-key", with: key

  # it turns out that using fill_in and triggering a .change() doesn't update Vue,
  # so let's jump through a few JS hoops instead
  # see https://github.com/vuejs/Discussion/issues/157
  page.execute_script(%{
    let elem = $("textarea[name='ssh-key']")[0];
    elem.value = '#{key}';
    elem.dispatchEvent(new Event('input', { 'bubbles': true }));
  })
end

Then %r{^my SSH Key should be added to the Existing Keys list$} do
  within(".existing-keys") do
    expect(page).to have_css("h3", text: "home-desktop")
  end
end

Then %r{^I should see my keys listed by name and fingerprint$} do
  within(".existing-keys") do
    expect(page).to have_css("li#ssh-public-key-1")
    expect(page).to have_css("h3", "home-desktop")
    expect(page).to have_content("SHA256:YwVZ0Zs7a3n6MiAK9jH6vrX8jbFDT0UwqWP76JQvlK4")
  end
end

Then %r{^each entry should have a 'Delete' button$} do
  within(".existing-keys li#ssh-public-key-1") do
    expect(page).to have_css("button", text: "Delete")
  end
end

When %r{^I click the 'Delete' button for my SSH key$} do
  within(".existing-keys li#ssh-public-key-1") do
    page.find("button", text: "Delete").click
  end
end

Then %r{^my key should have been deleted$} do
  expect(page).not_to have_css("#ssh-public-key-1")
end