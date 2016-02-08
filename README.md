# GoSquareToShopify
This project will convert a squarespace product list to a shopify one.

## Download Squarespace Data
You can find your squarespace data here: [https://your_domain_name/api/1/commerce/products] (GET).  This is easiest
by going to your site, and logging into the config page:  [https://your_domain_name/config] , logging in, then going
to the previous rest URL.  You should see lots of json data.  Save that off into a file.

## Convert SquareSpace data to Shopify

    $ GoSquareToShopify square.json [output.csv]

## Finally, upload your happy new CSV up to shopify
