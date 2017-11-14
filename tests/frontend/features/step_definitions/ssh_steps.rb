def user_with_key
  return {
    username: "rod.flanders",
    name: "Rod Flanders",
    email: "rod.flanders@springfield-elementary.k12.us",
    password: "okily-dokily!"
  }
end

Given %r{^my user account with public key exists$} do
  uri = URI('http://127.0.0.1:9095/setup/create_initial_user')
  req = Net::HTTP::Post.new(uri, "Content-Type" => "application/json")
  req.body = user_with_key.to_json
  res = Net::HTTP.start(uri.hostname, uri.port) do |http|
    http.request(req)
  end
  expect(res.code_type).to eql(Net::HTTPCreated)
end

Given %r{^my private key is valid$} do
  @ssh_key = "../backend/certificates/valid"
end

Given %r{^my private key is invalid$} do
  @ssh_key = "../backend/certificates/invalid"
end

When %r{^I initiate a SSH connection to the server$} do
  @response = connect_via_ssh("127.0.0.1", 2223, @ssh_key)
end

Then %r{^I should see the response "(.*?)"$} do |text|
  expect(@response).to include(text.strip)
end

Then %r{^I should receive an AuthenticationFailed error$} do
  expect(@response).to be_a(Net::SSH::AuthenticationFailed)
end

def connect_via_ssh(host, port, key, cmd="")

  response = StringIO.new

  Net::SSH.start(
    '127.0.0.1', 'git',
    port: port,
    host_key: "ssh-rsa",
    #encryption: "aes128-cbc",
    #compression: "zlib",
    keys: [key]
  ) do |session|

    session.open_channel do |channel|
      channel.on_data do |ch, data|
        response << "#{data}"
      end
      channel.exec ""
    end

    session.loop
  end

  return response.string

rescue Net::SSH::AuthenticationFailed => e
  return e
end