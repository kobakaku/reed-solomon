use reed_solomon::{Decoder, Encoder};

fn main() {
    let data = b"Hello World!";

    let ecc_len = 11;

    // Create encoder and decoder with
    let encoder = Encoder::new(ecc_len);
    let decorder = Decoder::new(ecc_len);

    // Encode data
    let encoded = encoder.encode(data);

    // Simulate some tranmission errors
    let mut corrupted = encoded;
    for i in 0..6 {
        corrupted[i] = 0x0;
    }

    // Try to recover data
    let known_erasures = [0];
    let recovered = decorder.correct(&corrupted, Some(&known_erasures)).unwrap();

    let orig_str = std::str::from_utf8(data).unwrap();
    let recv_str = std::str::from_utf8(recovered.data()).unwrap();

    println!("message:               {:?}", orig_str);
    println!("encoded:               {:?}", encoded);
    println!("corrupted:             {:?}", corrupted);
    println!("repaired:              {:?}", recv_str);
}
