Given %r{^I try to access "(.*?)" via HTTP$} do |path|
  uri = URI("http://127.0.0.1:9095/#{path}")
  req = Net::HTTP::Get.new(uri)
  @res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
end

Then %r{^I should receive a redirect to the same "(.*?)" but served via HTTPS$} do |path|
  expect(@res).to be_a(Net::HTTPMovedPermanently)
  expect(@res.body).to include("https://localhost:9096/#{path}")
end
