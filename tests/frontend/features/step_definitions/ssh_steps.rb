Given %r{^my user account with public key exists$} do


  # create the user
  user_uri = URI('http://127.0.0.1:9095/setup/create_initial_user')
  req = Net::HTTP::Post.new(user_uri, "Content-Type" => "application/json")
  req.body = user.to_json
  res = Net::HTTP.start(user_uri.hostname, user_uri.port) do |http|
    http.request(req)
  end
  expect(res.code_type).to eql(Net::HTTPCreated)

  # FIXME add steps for setting public key once complete
  step "I have logged in"
end

Given %r{^my private key is valid$} do
  step "I have an SSH key"
  @ssh_key = "../backend/certificates/valid"
end

Given %r{^my private key is invalid$} do
  @ssh_key = "../backend/certificates/invalid"
end

When %r{^I initiate a SSH connection to the server$} do
  @response = connect_via_ssh("127.0.0.1", 2223, @ssh_key)
end

Then %r{^I should see the response "(.*?)"$} do |text|
  expect(@response).to include(text)
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