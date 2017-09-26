Given %r{^I am on the appendix history page for "([^"]*)"$} do |arg1|
  path = "/cms/appendices/appendix_1.md/history"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I am on the document history page for "([^"]*)"$} do |arg1|
  path = "/cms/documents/document_1.md/history"
  visit(path)
  expect(page.current_path).to eql(path)
end
