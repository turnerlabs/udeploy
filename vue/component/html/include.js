export async function includeTenplates() {
    const tag = "html-template";

    let z, i, elmnt, file;

    z = document.getElementsByTagName("*");
    for (i = 0; i < z.length; i++) {
      elmnt = z[i];

      file = elmnt.getAttribute(tag);
      if (file) {
        const resp = await fetch(file);
        const html = await resp.text();

        if (resp.status == 200) {elmnt.innerHTML = html;}
        if (resp.status == 404) {elmnt.innerHTML = "Page not found.";}

        elmnt.removeAttribute(tag);

        return includeTenplates();
      }
    }
  }



