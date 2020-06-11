export default async function Fetcher(requestAddress, params, method = "GET", responseType="json", body = null){
    const url = new URL(requestAddress);
    url.search = new URLSearchParams(params).toString();
    console.log("url",url.toString());
    try{
        const response = await fetch(
            url.toString(),
            {method: method, body}
        );
        if (response.redirected){
            console.log("redirected");
            document.location = response.url;
            return [null, null]
        } else if (!response.ok){
            return [response.statusText, null];
        }else{
            return [null, await response[responseType]()]
        }
    }catch (err) {
        return [err.statusText, null];
    }
}