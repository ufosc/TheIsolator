require 'http'
require 'openssl'

document      = File.read('create-hello-world.json')
date          = Time.now.utc.httpdate
keypair       = OpenSSL::PKey::RSA.new(File.read('private.pem'))
signed_string = "(request-target): post /inbox\nhost: mastodon.social\ndate: #{date}"
signature     = Base64.strict_encode64(keypair.sign(OpenSSL::Digest::SHA256.new, signed_string))
header        = 'keyId="https://codecaptured.com/theisolator/activitypub/actor",headers="(request-target) host date",signature="' + signature + '"'

HTTP.headers({ 'Host': 'mastodon.social', 'Date': date, 'Signature': header })
    .post('https://mastodon.social/inbox', body: document)