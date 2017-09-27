Then %r{^I should only see the inactive '(.*?)' breadcrumb$} do |page_name|
  within("nav > ol.breadcrumb") do
    expect(page).to have_css("li", text: page_name)
  end
end

Then %r{^I should see the following breadcrumbs:$} do |table|
  actual = page.all("nav > ol.breadcrumb > li")
  expected = table.hashes.to_a

  fail "expected and actual sizes don't match" if actual.size != expected.size

  actual.zip(expected).each do |ac, ex|

    within(ac) do
      if ex['Reference'] == "None"
        expect(page).to have_content(ex['Text'])
      else
        expect(page).to have_link(ex['Text'], href: ex['Reference'])
      end
    end
  end

end
