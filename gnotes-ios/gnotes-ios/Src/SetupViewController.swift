//
//  SetupViewController.swift
//  gnotes-ios
//
//  Created by Westley Rose on 2/25/23.
//

import UIKit

class SetupViewController: UIViewController {

    @IBOutlet weak var fieldAccessKey: UITextField!
    @IBOutlet weak var fieldSecretKey: UITextField!
    @IBOutlet weak var fieldAccountID: UITextField!
    @IBOutlet weak var fieldCryptKey: UITextField!
    @IBOutlet weak var buttonDone: UIButton!


    override func viewDidLoad() {
        super.viewDidLoad()

        self.buttonDone.isEnabled = false;
    }


    @IBAction func buttonDone(_ sender: UIButton) {
        self.dismiss(animated: true)
    }

    @IBAction func fieldAccessKey(_ sender: UITextField) {
        self.view.endEditing(true)
        canEnableDoneButton()
    }
    @IBAction func fieldSecretKey(_ sender: UITextField) {
        self.view.endEditing(true)
        canEnableDoneButton()
    }
    @IBAction func fieldAccountID(_ sender: UITextField) {
        self.view.endEditing(true)
        canEnableDoneButton()
    }
    @IBAction func fieldCryptKey(_ sender: UITextField) {
        self.view.endEditing(true)
        canEnableDoneButton()
    }

    // MARK: - Navigation

    // In a storyboard-based application, you will often want to do a little preparation before navigation
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        // Get the new view controller using segue.destination.
        // Pass the selected object to the new view controller.
    }

    func canEnableDoneButton() {
        var entered = 0

        if self.fieldAccessKey.text != "" {
            entered += 1
        }
        if self.fieldSecretKey.text != "" {
            entered += 1
        }
        if self.fieldAccountID.text != "" {
            entered += 1
        }
        if self.fieldCryptKey.text != "" {
            entered += 1
        }

        if entered == 4 {
            self.buttonDone.isEnabled = true
        } else {
            self.buttonDone.isEnabled = false
        }
    }

}

