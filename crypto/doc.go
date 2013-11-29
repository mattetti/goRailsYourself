// Copyright (c) 2013, Matt Aimonetti
// Use of this source code is governed by a MIT style
// license that can be found at http://mattaimonetti.mit-license.org

/*

This crypto package ports some of Ruby on Rails' crypto (version 4+)
logic so encrypted/signed messages can be shared between a Ruby app and
a Go app. That said, this library is useful to anyone wanting to
encrypt/sign data.

The initial focus of this package was to be able to easily  share a Rails web session with a Go
app. Rails uses three classes provided by ActiveSupport (a library used
and maintained by the Rails team)
  * MessageEncryptor
  * MessageVerifier
  * KeyGenerator
to encrypt and sign sessions. In order to read/write a cookie session,
a Go app needs to be able to verify, decrypt/encrypt sign the session
data based on a shared secret.

Key components of this package

The main components of this package are:

  * MessageEncryptor
  * MessageVerifier 
  * KeyGenerator

The difference between MessageVerifier and MessageEncryptor is that you
want to use MessageEncryptor when you don't want the content of the data
to be available to people accessing the data. In both cases, the data is
signed but if the message is just signed, the content can be read.

Keygenerator is used to generate derived keys from a given secret.
If you want to generate a random key that isn't derived, look at
the GenerateRandomKey function.

Session serializer

As I'm writing this documentation, Rails doesn't let you change the
default session serializer (Marhsal). To be able to share the session
it needs to be serialized in a cross language format.
I wrote a patch to change Rails' default serializer to JSON:
https://gist.github.com/mattetti/7624413
Hopefully, an API will soon be available and JSON will become the default
serializer.
This package can use different serializers and you can also add your
own. This is useful if for instance you only have Go apps and choose to
use gob encoding or another encoding solution. Three serializers are
available JSON, XML and Null, the last serializer is basically a no-op
serializer used when the data doesn't need serialization and can be
transported as strings.

Rails session flow

It's important to understand how Rails handles the crypto around the
session.
Here is a quick and high level of what Rails does (Ruby code):

   # Secret set in the app.
    secret_key_base = "f7b5763636f4c1f3ff4bd444eacccca295d87b990cc104124017ad70550edcfd22b8e89465338254e0b608592a9aac29025440bfd9ce53579835ba06a86f85f9"
          
    key_generator = ActiveSupport::CachingKeyGenerator.new(ActiveSupport::KeyGenerator.new(secret_key_base, iterations: 1000))
    secret = key_generator.generate_key("encrypted cookie")
    sign_secret = key_generator.generate_key("signed encrypted cookie")
     
    encryptor = ActiveSupport::MessageEncryptor.new(secret, sign_secret, { serializer: JsonSessionSerializer } )
    # encrypt and sign the content of the session:
    encrypted_message = encryptor.encrypt_and_sign({msg: "hello world"})
    # The encrypted and signed message is stored in the session cookie
    # To decrypt and verify it:
    # encryptor.decrypt_and_verify(encrypted_message) # => {:msg => "hello world"}

The equivalent in Go is available in the documentation examples: http://godoc.org/github.com/mattetti/goRailsYourself/crypto#pkg-examples

Derived keys

A few important things need to be mentioned. Rails uses a unique secret
that is used to derive different keys using a default salt.
To read more about this process, see http://en.wikipedia.org/wiki/PBKDF2

Rails defaults to 1000 iterations when generating the derived keys, when
generating the keys in Go, we need to match the iteration number to get
the same keys. Note also that if the salt is changed in the Rails app,
you need to update it in your Go code.

There are two derived keys: one for encryption and one for signing. 
These keys are derived from the same secret but are different to
increase security. The keys are also of two different length.
The message signature is done by default using HMAC/sha1 requiring
a key of 64 bytes. However, the message is encrypted by default using
AES-256 CBC requiring a key of 32 bytes. 
Ruby's openssl library and this package automatically truncate longer
AES CBC keys so you can use two 64 byte keys. 
This is exactly what Rails does, it generates two keys of same length (64 bytes) and
lets the OpenSSL wrapper truncate the key. I, however recommend you
generate keys of different length to avoid any confusion. 
Here is an example:

  railsSecret := "f7b5763636f4c1f3ff4bd444eacccca295d87b990cc104124017ad70550edcfd22b8e89465338254e0b608592a9aac29025440bfd9ce53579835ba06a86f85f9"
  encryptedCookieSalt := []byte("encrypted cookie")
  encryptedSignedCookieSalt := []byte("signed encrypted cookie")

  kg := KeyGenerator{Secret: railsSecret}
  secret := kg.CacheGenerate(encryptedCookieSalt, 32)
  signSecret := kg.CacheGenerate(encryptedSignedCookieSalt, 64)
  e := MessageEncryptor{Key: secret, SignKey: signSecret}


Without Ruby

The encryption used in Rails isn't specific to Ruby and this library can
be used to communicate with apps that aren't in Ruby. As a matter of
fact, you might want to use this library to encrypt/sign your web
sessions/cookies even if you just have one Go app. The Rails implementation
has been tested and vested by many people and is safe to use.

*/
package crypto
